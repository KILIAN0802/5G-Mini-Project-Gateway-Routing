package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
	InstanceID     string `json:"instanceID"`
	ActiveRequests int    `json:"activeRequests"`
}

// Hàm xử lý
var sessions sync.Map

func CreateSession(
	w http.ResponseWriter,
	r *http.Request,
) {
	IncrementActiveRequest()
	defer DecrementActiveRequest()

	var req CreateSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if req.Supi != "" {
		sessions.Store(req.Supi, req)
	}

	w.Header().Set("Content-Type", "application/json")
	var buf [256]byte
	b := buf[:0]
	b = append(b, `{"handleby":"`...)
	b = append(b, instanceID...)
	b = append(b, `","status":"Active","pduSessionId":`...)
	b = strconv.AppendInt(b, int64(req.PduSessionId), 10)
	b = append(b, `,"supi":"`...)
	b = append(b, req.Supi...)
	b = append(b, `"}`...)
	w.Write(b)
}

func HealthCheck(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.Write([]byte("OK"))
}

func Metrics(
	w http.ResponseWriter,
	r *http.Request,
) {
	resp := MetricsResponse{
		InstanceID:     instanceID,
		ActiveRequests: int(atomic.LoadInt64(&activeRequests)),
	}

	w.Header().Set("Content-Type", "application/json")
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

	port := GetEnv("PORT", "9001")

	http.HandleFunc("/create-session", CreateSession)
	http.HandleFunc("/health", HealthCheck)

	http.HandleFunc(
		"/list-sessions",
		func(w http.ResponseWriter, r *http.Request) {
			parsedSessions := make(map[string]interface{})
			sessions.Range(func(key, value any) bool {
				if k, ok := key.(string); ok {
					parsedSessions[k] = value
				}
				return true
			})
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(parsedSessions)
		},
	)

	log.Println("PDU Session started: " + port)

	http.HandleFunc(
		"/metrics",
		Metrics,
	)

	h2s := &http2.Server{
		MaxConcurrentStreams:         50000,
		MaxReadFrameSize:             1048576,
		IdleTimeout:                  120 * time.Second,
		MaxUploadBufferPerStream:     65535 * 32,
		MaxUploadBufferPerConnection: 65535 * 64,
	}
	h2cHandler := h2c.NewHandler(http.DefaultServeMux, h2s)

	server := &http.Server{
		Addr:        ":" + port,
		Handler:     h2cHandler,
		IdleTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed:", err)
	}

}
