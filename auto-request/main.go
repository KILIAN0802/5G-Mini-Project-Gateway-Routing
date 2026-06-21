package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	targetURL      string
	TotalRequests  int
	MaxConcurrency int
	SupiBase       string
	ClientTimeout  time.Duration
}

var config = Config{
	targetURL:      "http://localhost:8080/nsmf-pdusession/v1/sm-contexts",
	TotalRequests:  5000,
	MaxConcurrency: 500,
	SupiBase:       "imsi-45204000000000",
	ClientTimeout:  30 * time.Second, // Tăng lên 30s để không bị timeout khi PDU Session xử lý mất 15s
}

// Metrics dùng atomic & mutex
var (
	successCount int64
	failCount    int64
	totalLatency int64
	handledMap   = make(map[string]int)
	mapMutex     sync.Mutex
)

type CreatSessionResponse struct {
	Handleby     string `json:"handleby"`
	Status       string `json:"status"`
	PduSessionId int    `json:"pduSessionId"`
	Supi         string `json:"supi"`
}

func resetMetrics() {
	atomic.StoreInt64(&successCount, 0)
	atomic.StoreInt64(&failCount, 0)
	atomic.StoreInt64(&totalLatency, 0)
	mapMutex.Lock()
	handledMap = make(map[string]int)
	mapMutex.Unlock()
}

func runBenchmark(client *http.Client, totalRequests int, maxConcurrency int) {
	resetMetrics()

	jobCh := make(chan int, maxConcurrency)
	var wg sync.WaitGroup

	for i := 0; i < maxConcurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for reqID := range jobCh {
				sendRequest(client, reqID)
			}
		}(i)
	}

	start := time.Now()
	for i := 0; i < totalRequests; i++ {
		jobCh <- i
	}

	close(jobCh)
	wg.Wait()

	duration := time.Since(start)

	log.Printf("=== RESULT ===")
	log.Printf("Duration: %s", duration)
	log.Printf("Total Requests: %d", totalRequests)
	log.Printf("Success: %d", atomic.LoadInt64(&successCount))
	log.Printf("Failed: %d", atomic.LoadInt64(&failCount))
	
	success := atomic.LoadInt64(&successCount)
	if success > 0 {
		log.Printf("Avg Latency: %.2fms", float64(atomic.LoadInt64(&totalLatency))/float64(success))
	} else {
		log.Printf("Avg Latency: 0.00ms")
	}
	log.Printf("Error Percent: %.2f%%", float64(atomic.LoadInt64(&failCount))/float64(totalRequests)*100)

	log.Printf("=== LOAD DISTRIBUTION ===")
	mapMutex.Lock()
	for instance, count := range handledMap {
		var pct float64
		if success > 0 {
			pct = float64(count) / float64(success) * 100
		}
		log.Printf("Instance %s handled: %d requests (%.2f%%)", instance, count, pct)
	}
	mapMutex.Unlock()
}

func main() {
	// Khởi tạo bộ sinh số ngẫu nhiên
	rand.Seed(time.Now().UnixNano())

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			MaxIdleConns:        10000,
			MaxConnsPerHost:     1000,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     30 * time.Second,
		},
		Timeout: config.ClientTimeout,
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== CẤU HÌNH KIỂM THỬ ===")
	fmt.Printf("Đường dẫn đích (Target): %s\n", config.targetURL)
	fmt.Printf("HTTP Client Timeout: %v\n", config.ClientTimeout)
	fmt.Println("=========================")
	fmt.Println("Chọn chế độ bắn request:")
	fmt.Println("1. Tự động bắn theo chu kỳ 5 - 7 giây (mỗi chu kỳ 500 requests, nhấn ESC để dừng)")
	fmt.Println("2. Nhập thủ công số lượng từ terminal (nhấn ESC hoặc gõ 'exit' để dừng)")
	fmt.Print("Lựa chọn của bạn (1 hoặc 2): ")

	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)

	if choiceStr == "1" {
		fmt.Println("\n>>> ĐANG CHẠY CHẾ ĐỘ 1: Bắn theo chu kỳ 5 - 7 giây. Nhấn phím ESC để dừng...")
		cycle := 1
		for {
			if IsEscPressed() {
				fmt.Println("\n[!] Đã nhận phím ESC. Đang dừng chương trình...")
				break
			}

			fmt.Printf("\n--- [Chu kỳ %d] Bắt đầu bắn 500 requests ---\n", cycle)
			runBenchmark(client, 500, 500)

			// Ngẫu nhiên chu kỳ 5 đến 7 giây
			intervalSeconds := 5 + rand.Intn(3)
			fmt.Printf("--- [Chu kỳ %d] Hoàn thành. Đợi %d giây cho chu kỳ tiếp theo (Hoặc nhấn ESC để dừng)... ---\n", cycle, intervalSeconds)

			// Chia nhỏ giấc ngủ để liên tục kiểm tra phím ESC
			stop := false
			sleepSteps := intervalSeconds * 10
			for i := 0; i < sleepSteps; i++ {
				if IsEscPressed() {
					fmt.Println("\n[!] Đã nhận phím ESC. Đang dừng chương trình...")
					stop = true
					break
				}
				time.Sleep(100 * time.Millisecond)
			}
			if stop {
				break
			}
			cycle++
		}
	} else if choiceStr == "2" {
		fmt.Println("\n>>> ĐANG CHẠY CHẾ ĐỘ 2: Nhập thủ công số lượng từ Terminal. Nhấn phím ESC hoặc gõ 'exit' để dừng...")
		for {
			if IsEscPressed() {
				fmt.Println("\n[!] Đã nhận phím ESC. Đang dừng chương trình...")
				break
			}

			fmt.Print("\nNhập số lượng request muốn bắn: ")
			
			// Đọc input một cách phi chặn để cho phép nhấn ESC dừng chương trình khi đang đợi input
			inputCh := make(chan string, 1)
			go func() {
				input, _ := reader.ReadString('\n')
				inputCh <- input
			}()

			var input string
			stop := false
			for {
				select {
				case input = <-inputCh:
					// Nhận dữ liệu nhập từ bàn phím
				case <-time.After(100 * time.Millisecond):
					// Liên tục kiểm tra phím ESC
					if IsEscPressed() {
						fmt.Println("\n[!] Đã nhận phím ESC. Đang dừng chương trình...")
						stop = true
					}
				}
				if input != "" || stop {
					break
				}
			}

			if stop {
				break
			}

			input = strings.TrimSpace(input)
			if input == "exit" || input == "quit" {
				fmt.Println("Đang dừng chương trình...")
				break
			}

			numRequests, err := strconv.Atoi(input)
			if err != nil || numRequests <= 0 {
				fmt.Println("Số lượng không hợp lệ. Vui lòng nhập số nguyên dương.")
				continue
			}

			concurrency := numRequests
			if concurrency > 500 {
				concurrency = 500 // Giới hạn concurrency tối đa là 500 để đảm bảo tài nguyên
			}

			fmt.Printf("Đang bắn đồng thời %d requests (Concurrency: %d)...\n", numRequests, concurrency)
			runBenchmark(client, numRequests, concurrency)
		}
	} else {
		fmt.Println("Lựa chọn không hợp lệ. Chương trình kết thúc.")
	}
}

// sendRequest gửi một request và tính toán độ trễ, số liệu
func sendRequest(client *http.Client, reqID int) {
	requestData := map[string]interface{}{
		"supi":         fmt.Sprintf("%s%06d", config.SupiBase, reqID),
		"gpsi":         "0919213419",
		"pduSessionId": reqID,
		"dnn":          "internet",
		"sNssai": map[string]interface{}{
			"sst": 1,
			"sd":  "000001",
		},
		"servingNfid": "amf-1",
		"anType":       "3GPP",
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