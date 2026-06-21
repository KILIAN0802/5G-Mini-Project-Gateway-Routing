package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

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
	InstanceID   string `json:"instanceID"`
	ActiveRequests int    `json:"activeRequests"`
}

// Hàm xử lý
var mu sync.Mutex

func CreateSession(
	w http.ResponseWriter, //
	r *http.Request,
) {
	IncrementActiveRequest()
	defer DecrementActiveRequest()
	
	delayMode := GetEnv("DELAY_MODE", "fixed")
	var delayDuration time.Duration

	if delayMode == "random" {
		// Sinh ngẫu nhiên thời gian xử lý: random % 20 giây (0 -> 19s)
		delaySeconds := rand.Intn(20)
		delayDuration = time.Duration(delaySeconds) * time.Second
	} else {
		// Mặc định cố định 15 giây
		delayDuration = 15 * time.Second
	}

	log.Printf("[%s] Bat dau xu ly session, sleep %v", instanceID, delayDuration)
	time.Sleep(delayDuration)
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
	reqJSON, errSave := json.Marshal(req)
	if errSave == nil {
		errRedis := SaveSessionInRedis(req.Supi, string(reqJSON))
		if errRedis != nil {
			log.Printf("[%s] Lưu session vào Redis thất bại: %v", instanceID, errRedis)
		}
	}else{
		log.Printf("[%s] Chuyển request thành JSON thất bại: %v", instanceID, errSave)
	}
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
	resp := MetricsResponse{
		InstanceID: instanceID,
		ActiveRequests: int(atomic.LoadInt64(&activeRequests)),
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	json.NewEncoder(w).Encode(resp)
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
	
	http.ListenAndServe(
		":"+port, // Cổng
		nil,
	)

}
