package main

import (
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"github.com/shirou/gopsutil/v3/cpu"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func getCGroupCPUUsage() (int64, error) {
	// 1. Thử đọc cgroup v2 (/sys/fs/cgroup/cpu.stat)
	data, err := os.ReadFile("/sys/fs/cgroup/cpu.stat")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "usage_usec") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					usec, err := strconv.ParseInt(fields[1], 10, 64)
					if err == nil {
						return usec * 1000, nil
					}
				}
			}
		}
	}

	// 2. Thử đọc cgroup v1 (/sys/fs/cgroup/cpu/cpuacct.usage hoặc /sys/fs/cgroup/cpuacct/cpuacct.usage)
	paths := []string{
		"/sys/fs/cgroup/cpu/cpuacct.usage",
		"/sys/fs/cgroup/cpuacct/cpuacct.usage",
	}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil {
			ns, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
			if err == nil {
				return ns, nil
			}
		}
	}

	return 0, fmt.Errorf("cgroup cpu metrics not found")
}

type CGroupCPUSampler struct {
	lastUsage int64
	lastTime  time.Time
}

func (s *CGroupCPUSampler) SamplePercent() float64 {
	now := time.Now()
	ns, err := getCGroupCPUUsage()
	if err != nil {
		percentages, err := cpu.Percent(0, false)
		if err == nil && len(percentages) > 0 {
			return percentages[0]
		}
		return 0
	}

	if s.lastUsage == 0 || s.lastTime.IsZero() {
		s.lastUsage = ns
		s.lastTime = now
		return 0
	}

	deltaUsage := float64(ns - s.lastUsage)
	deltaTime := float64(now.Sub(s.lastTime).Nanoseconds())

	s.lastUsage = ns
	s.lastTime = now

	if deltaTime <= 0 {
		return 0
	}

	pct := (deltaUsage / deltaTime) * 100.0
	if pct < 0 {
		pct = 0
	}
	return pct
}

type SNssai struct {
	SST int    `json:"sst"` // Single-Use Scenario ID - Xác định loại hình dịch vụ hoặc tập hợp các tính năng.
	SD  string `json:"sd"`  // Single-Use Scenario ID - Xác định loại hình dịch vụ hoặc tập hợp các tính năng.
}

type CreateSessionRequest struct {
	Supi         string `json:"supi"` // Subscription Permanent Identifier - Lưu ở USIM/ eSIM và trong UDM
	Gpsi         string `json:"gpsi"` // Generetic Public Subscription Identifier - Số điện thoại
	PduSessionId int    `json:"pduSessionId"`
	Dnn          string `json:"dnn"` // Data Network Name - VD: internet, IMS, mạng riêng doanh nghiệp, ...
	SNssai       SNssai `json:"sNssai"`
	ServingNfid  string `json:"servingNfid"` // Serving Network Function Identifier
	AnType       string `json:"anType"`      // Access Type - Loại kết nối (vd: 3gpp, non-3gpp)
}

type CreateSessionResponse struct {
	Handleby     string `json:"handleby"`
	Status       string `json:"status"`
	PduSessionId int    `json:"pduSessionId"`
	Supi         string `json:"supi"`
}

type MetricsResponse struct {
	InstanceID     string  `json:"instanceID"`
	ActiveRequests int     `json:"activeRequests"`
	PeakCPU        float64 `json:"peakCpu"`
}

// Hàm xử lý
var mu sync.Mutex
var (
	peakCPU      float64
	peakCPUMutex sync.Mutex
)

func CreateSession(
	w http.ResponseWriter, //
	r *http.Request,
) {
	IncrementActiveRequest()
	defer DecrementActiveRequest()
	
	// delayMode := GetEnv("DELAY_MODE", "fixed")
	// var delayDuration time.Duration

	// if delayMode == "random" {
	// 	// Sinh ngẫu nhiên thời gian xử lý: random % 20 giây (0 -> 19s)
	// 	delaySeconds := rand.Intn(20)
	// 	delayDuration = time.Duration(delaySeconds) * time.Second
	// } else {
	// 	// Mặc định cố định 15 giây
	// 	delayDuration = 15 * time.Second
	// }

	// log.Printf("[%s] Bat dau xu ly session, sleep %v", instanceID, delayDuration)
	// time.Sleep(delayDuration)
	var req CreateSessionRequest

	err := json.NewDecoder(
		r.Body,
	).Decode(&req)

	if err != nil {
		http.Error(
			w,
			"bad request",         // Response message -> Ghi vào r.Body
			400, // 400 = Bad Request
		)
		return // Trả về r.Body và dừng hàm
	}

	resp := CreateSessionResponse{
		Handleby:     instanceID,
		Status:       "Active",
		PduSessionId: req.PduSessionId,
		Supi:         req.Supi,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	) // Thiết lập HTTP header cho response
	// w.WriteHeader(200)// 200 OK
	// reqJSON, errSave := json.Marshal(req)
	// if errSave == nil {
	// go func ()  {
	// 	errRedis := SaveSessionInRedis(req.Supi, string(reqJSON))
	// 	if errRedis != nil {
	// 		log.Printf("[%s] Lưu session vào Redis thất bại: %v", instanceID, errRedis)
	// 	}
	// }()
	// }else{
	// 	log.Printf("[%s] Chuyển request thành JSON thất bại: %v", instanceID, errSave)
	// }
	json.NewEncoder(w).Encode(resp)

}



func HealthCheck(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Write(
		[]byte("OK"), // Chuyển chuỗi "OK" thành mảng byte vì hàm Write yêu cầu mảng byte
	)
}

func Metrics(
	w http.ResponseWriter,
	r *http.Request,
){
	peakCPUMutex.Lock()
	pCPU := peakCPU
	peakCPUMutex.Unlock()

	resp := MetricsResponse{
		InstanceID:     instanceID,
		ActiveRequests: int(atomic.LoadInt64(&activeRequests)),
		PeakCPU:        pCPU,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	json.NewEncoder(w).Encode(resp)
}

func ResetMetrics(
	w http.ResponseWriter,
	r *http.Request,
) {
	peakCPUMutex.Lock()
	peakCPU = 0
	peakCPUMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

var instanceID string

var activeRequests int64

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown-ip"
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown-ip"
}

func main() {
	rand.Seed(time.Now().UnixNano())
	instanceID = GetEnv("INSTANCE_ID", "")
	initRedis()
	if instanceID == "" {
		host, err := os.Hostname()
		if err != nil {
			host = "pdu-unknown"
		}
		ip := getLocalIP()
		instanceID = host + " (" + ip + ")"
	}

	// Goroutine ngầm đo % CPU của container pdu-session theo cửa sổ 1s (khớp 100% với docker stats)
	go func() {
		sampler := &CGroupCPUSampler{}
		ticker := time.NewTicker(1000 * time.Millisecond)
		for range ticker.C {
			pct := sampler.SamplePercent()
			if pct > 0 {
				peakCPUMutex.Lock()
				if pct > peakCPU {
					peakCPU = pct
				}
				peakCPUMutex.Unlock()
			}
		}
	}()

	port := GetEnv(
		"PORT",
		"9001",
	)

	http.HandleFunc(
		"/create-session",
		CreateSession,
	)
	http.HandleFunc(
		"/health",
		HealthCheck,
	)

	http.HandleFunc(
		"/list-sessions",
		func(w http.ResponseWriter, r *http.Request){
			log.Printf("[%s] API list-sessions duoc goi", instanceID)
			data, err := GetAllSessionsFromRedis()
			if err != nil {
				log.Printf("[%s] Loi doc tu Redis: %v", instanceID, err)
				http.Error(w, "Redis read error: "+err.Error(), 500)
				return
			}
			parsedSessions := make(map[string] interface{})
			for supi, rawJSON := range data {
				var val interface{}
				if err := json.Unmarshal([]byte(rawJSON), &val); err == nil {
					parsedSessions[supi] = val
				}else {
					parsedSessions[supi] = rawJSON
				}
			}
			w.Header().Set(
				"Content-Type",
				"application/json",
			)
			json.NewEncoder(w).Encode(parsedSessions)
		},
	)

	log.Println("PDU Session started: " + port)

	http.HandleFunc(
		"/metrics",
		Metrics,
	)
	http.HandleFunc(
		"/reset-metrics",
		ResetMetrics,
	)

	h2s := &http2.Server{}
	h2cHandler := h2c.NewHandler(http.DefaultServeMux, h2s)

	server := &http.Server{
		Addr: ":" + port,
		Handler: h2cHandler,
		IdleTimeout: 10 * time.Second, 
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed:", err)
	}

}
