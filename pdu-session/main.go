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
