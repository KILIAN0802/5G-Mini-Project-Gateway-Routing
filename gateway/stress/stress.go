package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// Số lượng request tổng cộng
	totalRequests := 25
	// Mỗi request cách nhau 1 giây để người dùng dễ quan sát tải tăng/giảm dần dần
	spawnInterval := 1000 * time.Millisecond
	
	fmt.Printf("Bắt đầu test từ từ với %d requests (mỗi %v gửi 1 request)...\n", totalRequests, spawnInterval)

	var wg sync.WaitGroup

	// Goroutine in trạng thái tải mỗi 1 giây
	stopMonitor := make(chan struct{})
	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				resp, err := http.Get("http://localhost:8080/instances")
				if err != nil {
					continue
				}
				var debugData []map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&debugData)
				resp.Body.Close()
				
				fmt.Printf("[%s] Số request đang xử lý: ", time.Now().Format("15:04:05"))
				for _, inst := range debugData {
					fmt.Printf("%s: %v | ", inst["id"], inst["activeRequests"])
				}
				fmt.Println()
			case <-stopMonitor:
				return
			}
		}
	}()

	startTime := time.Now()
	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			body := []byte(`{"supi":"gradual-user","pduSessionId":1}`)
			resp, err := http.Post("http://localhost:8080/nsmf-pdusession/v1/sm-contexts", "application/json", bytes.NewBuffer(body))
			if err != nil {
				return
			}
			resp.Body.Close()
		}(i)
		
		// Đợi 1 khoảng thời gian trước khi gửi request tiếp theo
		time.Sleep(spawnInterval)
	}

	wg.Wait()
	close(stopMonitor)
	fmt.Printf("Hoàn thành test trong %v!\n", time.Since(startTime))
}
