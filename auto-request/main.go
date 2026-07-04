package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
)

type Config struct {
	targetURL     string
	TotalRequests int
	SupiBase      string
	ClientTimeout time.Duration
}

var config = Config{
	targetURL:     "http://127.0.0.1:8080/nsmf-pdusession/v1/sm-contexts",
	TotalRequests: 5000,
	SupiBase:      "imsi-45204000000000",
	ClientTimeout: 90 * time.Second, // Tăng lên 30s để không bị timeout khi PDU Session xử lý mất 15s
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

func runBenchmark(clients []*http.Client, TotalRequests int, concurrency int) {
	resetMetrics()

	var wg sync.WaitGroup
	jobs := make(chan int, TotalRequests)

	// Đẩy toàn bộ request ID vào queue jobs
	for i := 0; i < TotalRequests; i++ {
		jobs <- i
	}
	close(jobs)

	// Nếu concurrency <= 0, mặc định bằng TotalRequests (bắn đồng thời tối đa)
	if concurrency <= 0 {
		concurrency = TotalRequests
	}

	startTime := time.Now()

	// Khởi chạy worker pool
	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			client := clients[workerID % len(clients)]
			for reqID := range jobs {
				requestData := map[string]interface{}{
					"supi":         fmt.Sprintf("%s%010d", config.SupiBase, reqID),
					"gspi":         "0919213638",
					"pduSessionId": reqID,
					"dnn":          "internet",
					"sNssai": map[string]interface{}{
						"sst": 1,
						"sd":  "000001",
					},
					"servingNfid": "amf-1",
					"anType":      "3GPP",
				}

				jsonData, err := json.Marshal(requestData)
				if err != nil {
					atomic.AddInt64(&failCount, 1)
					continue
				}

				reqStartTime := time.Now()
				resp, err := client.Post(config.targetURL, "application/json", bytes.NewBuffer(jsonData))

				if err != nil {
					atomic.AddInt64(&failCount, 1)
					log.Printf("[ERR] Request %d failed to send: %v", reqID, err)
					continue
				}

				latency := time.Since(reqStartTime).Milliseconds()

				if resp.StatusCode != 200 {
					atomic.AddInt64(&failCount, 1)
					bodyBytes, _ := io.ReadAll(resp.Body)
					log.Printf("[ERR] Request %d received status %d: %s", reqID, resp.StatusCode, string(bodyBytes))
					resp.Body.Close()
					continue
				}

				var result CreatSessionResponse
				err = json.NewDecoder(resp.Body).Decode(&result)
				resp.Body.Close()
				if err != nil {
					atomic.AddInt64(&failCount, 1)
					log.Println(err)
					continue
				}

				expectedSupi := fmt.Sprintf("%s%010d", config.SupiBase, reqID)
				if result.PduSessionId != reqID || result.Supi != expectedSupi {
					log.Printf("[LỖI] Sai luồng! Gửi ID: %d (Nhận ID: %d) | Gửi SUPI: %s (Nhận SUPI: %s)", 
						reqID, result.PduSessionId, expectedSupi, result.Supi)
					atomic.AddInt64(&failCount, 1)
					continue
				}

				atomic.AddInt64(&successCount, 1)
				atomic.AddInt64(&totalLatency, latency)
				mapMutex.Lock()
				handledMap[result.Handleby]++
				mapMutex.Unlock()
			}
		}(w)
	}

	wg.Wait()

	duration := time.Since(startTime)
	tps := float64(TotalRequests) / duration.Seconds()
	log.Printf("=== KẾT QUẢ TEST TẢI ĐỒNG THỜI ===")
	log.Printf("Thời gian chạy: %s", duration)
	log.Printf("Tổng số Request gửi đi: %d", TotalRequests)
	log.Printf("Thành công: %d", atomic.LoadInt64(&successCount))
	log.Printf("Thất bại (lỗi kết nối/timeout/sai luồng): %d", atomic.LoadInt64(&failCount))
	log.Printf("TPS: %.2f", tps)

	success := atomic.LoadInt64(&successCount)
	// if success > 0 {
	// 	log.Printf("Độ trễ trung bình: %.2fms", float64(atomic.LoadInt64(&totalLatency))/float64(success))
	// }

	log.Printf("=== PHÂN PHỐI TẢI TRÊN CÁC INSTANCES ===")
	mapMutex.Lock()
	for instance, count := range handledMap {
		var pct float64
		if success > 0 {
			pct = float64(count) / float64(success) * 100
		}
		log.Printf("Instance %s xử lý: %d requests (%.2f%%)", instance, count, pct)
	}
	mapMutex.Unlock()
}


func main() {
	// Khởi tạo bộ sinh số ngẫu nhiên
	rand.Seed(time.Now().UnixNano())

	if envURL := os.Getenv("TARGET_URL"); envURL != "" {
		config.targetURL = envURL
	}

	//Notice
	client := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, network, addr)
			},
		},
		Timeout: config.ClientTimeout,
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== CẤU HÌNH KIỂM THỬ ===")
	fmt.Printf("Đường dẫn đích (Target): %s\n", config.targetURL)
	fmt.Printf("HTTP Client Timeout: %v\n", config.ClientTimeout)
	fmt.Println("Chế độ: Bắn đồng thời toàn bộ luồng (Không giới hạn)")
	fmt.Println("=========================")
	fmt.Println("Chọn chế độ bắn request:")
	fmt.Println("1. Tự động bắn theo chu kỳ 5 - 7 giây (mỗi chu kỳ 500 requests, nhấn ESC để dừng)")
	fmt.Println("2. Nhập thủ công số lượng từ terminal (Thiết lập số request, connection và worker)")
	fmt.Print("Lựa chọn của bạn (1 hoặc 2): ")

	choiceStr, _ := reader.ReadString('\n')
	choiceStr = strings.TrimSpace(choiceStr)

	if choiceStr == "1" {
		cycle := 1
		for {
			if IsEscPressed() {
				break
			}

			fmt.Printf("\n--- [Chu kỳ %d] Bắt đầu bắn 500 requests ---\n", cycle)
			runBenchmark([]*http.Client{client}, 500, 500)

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
		for {
			if IsEscPressed() {
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
				break
			}

			numRequests, err := strconv.Atoi(input)
			if err != nil || numRequests <= 0 {
				fmt.Println("Số lượng không hợp lệ. Vui lòng nhập số nguyên dương.")
				continue
			}

			// Hỏi số lượng connection
			fmt.Print("Nhập số lượng connection (mặc định 1): ")
			inputConnCh := make(chan string, 1)
			go func() {
				inputC, _ := reader.ReadString('\n')
				inputConnCh <- inputC
			}()

			var inputConn string
			for {
				select {
				case inputConn = <-inputConnCh:
				case <-time.After(100 * time.Millisecond):
					if IsEscPressed() {
						stop = true
					}
				}
				if inputConn != "" || stop {
					break
				}
			}

			if stop {
				break
			}

			inputConn = strings.TrimSpace(inputConn)
			numConnections, err := strconv.Atoi(inputConn)
			if err != nil || numConnections <= 0 {
				numConnections = 1
			}

			// Hỏi số lượng worker song song
			fmt.Print("Nhập số lượng Worker song song (mặc định 100): ")
			inputWorkerCh := make(chan string, 1)
			go func() {
				inputW, _ := reader.ReadString('\n')
				inputWorkerCh <- inputW
			}()

			var inputWorker string
			for {
				select {
				case inputWorker = <-inputWorkerCh:
				case <-time.After(100 * time.Millisecond):
					if IsEscPressed() {
						stop = true
					}
				}
				if inputWorker != "" || stop {
					break
				}
			}

			if stop {
				break
			}

			inputWorker = strings.TrimSpace(inputWorker)
			numWorkers, err := strconv.Atoi(inputWorker)
			if err != nil || numWorkers <= 0 {
				numWorkers = 100
			}

			// Tạo client pool tương ứng với số connection mong muốn
			clients := make([]*http.Client, numConnections)
			for i := 0; i < numConnections; i++ {
				clients[i] = &http.Client{
					Transport: &http2.Transport{
						AllowHTTP: true,
						DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
							var d net.Dialer
							return d.DialContext(ctx, network, addr)
						},
					},
					Timeout: config.ClientTimeout,
				}
			}

			fmt.Printf("Đang bắn %d requests qua %d TCP Connections song song với %d Workers\n", numRequests, numConnections, numWorkers)
			runBenchmark(clients, numRequests, numWorkers)
		}
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
		"anType":      "3GPP",
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
