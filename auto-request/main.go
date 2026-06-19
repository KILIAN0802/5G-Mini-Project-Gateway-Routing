package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"sync"
	"sync/atomic"
)


type Config struct {
	targetURL string
	TotalRequests int
	MaxConcurrency int
	SupiBase string
	ClientTimeout time.Duration
}

var config = Config{
	targetURL: "http://localhost:8080/nsmf-pdusession/v1/sm-contexts",
	TotalRequests: 5000,
	MaxConcurrency: 500,
	SupiBase: "imsi-45204000000000",
	ClientTimeout: 10 * time.Second,
}

// Metrics dùng atomic & mutex
var (
	successCount int64
	failCount int64
	totalLatency int64
	handledMap = make(map[string]int)
	mapMutex sync.Mutex
)
type CreatSessionResponse struct {
    Handleby     string `json:"handleby"`
    Status       string `json:"status"`
    PduSessionId int    `json:"pduSessionId"`
    Supi         string `json:"supi"`
}



func main() {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns: 10000,
			MaxConnsPerHost: 1000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout: 30 * time.Second,
		},
		Timeout: config.ClientTimeout,
	}

	jobCh := make(chan int, config.MaxConcurrency)
	var wg sync.WaitGroup

	for i := 0; i < config.MaxConcurrency; i++ {
		wg.Add(1) //tăng wg lên 1
		go func(workerID int){
			defer wg.Done() // Khi công việc hoàn thành thì wg giảm 1 
			for reqID := range jobCh {
				sendRequest(client, reqID)
			}
		}(i)
	}

	start := time.Now()
	for i := 0; i < config.TotalRequests; i++ {
		jobCh <- i
	}

	close(jobCh)
	wg.Wait()

	duration := time.Since(start)
	
	log.Printf("=== RESULT ===")
	log.Printf("Duration: %s", duration)
	log.Printf("Total Requests: %d", config.TotalRequests)
	log.Printf("Success: %d", atomic.LoadInt64(&successCount))
	log.Printf("Failed: %d", atomic.LoadInt64(&failCount))
	log.Printf("Avg Latency: %.2fms", float64(atomic.LoadInt64(&totalLatency)) / float64(atomic.LoadInt64(&successCount)))
	log.Printf("Error Percent: %.2f%%", float64(atomic.LoadInt64(&failCount)) / float64(config.TotalRequests) * 100)

	log.Printf("=== LOAD DISTRIBUTION ===")
	mapMutex.Lock()
	for instance, count := range handledMap {
		log.Printf("Instance %s handled: %d requests (%.2f%%)", 
			instance, count, float64(count)/float64(atomic.LoadInt64(&successCount))*100)
	}
	mapMutex.Unlock()
}	

// Do không nằm chung thư mục với pdu-session
func sendRequest(client *http.Client, reqID int) {
	
	requestData := map[string]interface{}{
		"supi": fmt.Sprintf("%s%06d", config.SupiBase, reqID),
		"gpsi": "0919213419",
		"pduSessionId": reqID,
		"dnn": "internet",
		"sNssai": map[string]interface{}{
			"sst": 1,
			"sd": "000001",
		},
		"servingNfid": "amf-1",
		"anType": "3GPP",
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		atomic.AddInt64(&failCount, 1)
		return
	}

	startTime := time.Now()
	resp, err := client.Post(config.targetURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		atomic.AddInt64(&failCount, 1)
		return
	}
	defer resp.Body.Close()

	latency := time.Since(startTime).Milliseconds()
	
	if resp.StatusCode != 200 {
		atomic.AddInt64(&failCount, 1)
		return 
	}


	var result CreatSessionResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		atomic.AddInt64(&failCount, 1)
		return
	}



	atomic.AddInt64(&successCount, 1)
	atomic.AddInt64(&totalLatency, latency)

	mapMutex.Lock()
	handledMap[result.Handleby]++
	mapMutex.Unlock()
}