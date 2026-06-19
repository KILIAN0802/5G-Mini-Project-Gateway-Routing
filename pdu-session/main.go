package main

import (
	"encoding/json"
	"os"
	// "fmt"
	"log"
	"net/http"
	"sync"
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
	
	time.Sleep(50 * time.Millisecond)
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
		ActiveRequests: GetActiveRequests(),
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)
	json.NewEncoder(w).Encode(resp)
}

var instanceID string

var activeRequests int


func main() {
	instanceID = GetEnv("INSTANCE_ID", "")
	if instanceID == "" {
		if host, err := os.Hostname(); err == nil {
			instanceID = host
		}else{
			instanceID = "pdu-unknown"
		}
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
