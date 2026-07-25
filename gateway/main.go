package main

import (
	// "bytes" // Cung cấp các hàm để thao tác với kiểu dữ liệu byte -> Thường dùng để tạo buffer cho việc đọc / ghi dữ liệu
	"context"
	"crypto/tls"
	// "encoding/json"
	"io" // Cung cấp các interface chuẩn để đọc/ghi dữ liệu
	"log"
	"net"
	"net/http"
	"sync"
	// "sync/atomic"

	// Cho phép:
	// Tạo web server ( http.ListenAndServe)
	// Gửi request HTTP (http.Get, http.Post)
	// Xử lý request/response qua http.Handler
	"gateway/algorithm"
	"gateway/handler"
	"gateway/health"
	"gateway/registry"
	// "strconv"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 4096)
		return &b
	},
}

var pduClient = &http.Client{
	Timeout: 90 * time.Second,
	Transport: &http2.Transport{
		AllowHTTP: true,
		DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 60 * time.Second,
			}
			conn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			if tcpConn, ok := conn.(*net.TCPConn); ok {
				_ = tcpConn.SetNoDelay(true)
			}
			return conn, nil
		},
		StrictMaxConcurrentStreams: false,
		ReadIdleTimeout:            30 * time.Second,
		PingTimeout:                15 * time.Second,
	},
}

func ForwardToPDU(
	w http.ResponseWriter,
	r *http.Request,
) {
	selected := algorithm.SelectBackend(registry.GetHealthyInstance())
	if selected == nil {
		http.Error(w, "NO_BACKEND_AVAILABLE", 500)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), "POST", "http://"+selected.Address+"/create-session", r.Body)
	if err != nil {
		http.Error(w, "Error creating request", 500)
		return
	}
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))

	resp, err := pduClient.Do(req)
	if err != nil {
		log.Printf("ForwardToPDU error forwarding to %s: %v", selected.Address, err)
		http.Error(w, "Backend Error: "+err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func ListSessionsForward(w http.ResponseWriter, r *http.Request) {
	selected := algorithm.SelectBackend(registry.GetHealthyInstance())
	if selected == nil {
		http.Error(w, "NO_BACKEND_AVAILABLE", 500)
		return
	}

	req, err := http.NewRequest("GET", "http://"+selected.Address+"/list-sessions", nil)
	if err != nil {
		http.Error(w, "Error creating request", 500)
		return
	}

	resp, err := pduClient.Do(req)
	if err != nil {
		log.Printf("ListSessionsForward error forwarding to %s: %v", selected.Address, err)
		http.Error(w, "Backend Error: "+err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

// RateLimiter triển khai thuật toán Token Bucket (Xô chứa Token) để kiểm soát tốc độ request (RPS)
type RateLimiter struct {
	rate       float64    // Tốc độ hồi (nạp) token mới vào xô mỗi giây (Số token/giây)
	capacity   float64    // Dung tích tối đa của xô (Giới hạn tối đa số token tích lũy / Ngưỡng bùng nổ Burst)
	tokens     float64    // Số lượng token hiện có sẵn trong xô tại thời điểm kiểm tra
	lastRefill time.Time  // Mốc thời gian của lần tính toán / nạp token gần đây nhất
	mu         sync.Mutex // Mutex khóa đồng bộ để đảm bảo an toàn đa luồng (Thread-safe) khi nhiều Goroutines gọi đồng thời
}

// NewRateLimiter khởi tạo một bộ giới hạn tốc độ mới với tốc độ nạp rate và dung tích xô capacity
func NewRateLimiter(rate float64, capacity float64) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity, // Ban đầu xô được nạp đầy token
		lastRefill: time.Now(),
	}
}

// Allow kiểm tra xem request hiện tại có được phép đi tiếp hay không
func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	// Tính thời gian đã trôi qua kể từ lần nạp token gần nhất (tính bằng giây)
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.lastRefill = now

	// Nạp thêm token tương ứng với khoảng thời gian đã trôi qua: (thời gian trôi qua * tốc độ rate)
	rl.tokens += elapsed * rl.rate
	// Không cho phép số token vượt quá dung tích tối đa của xô (capacity)
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}

	// Nếu trong xô còn đủ ít nhất 1 token -> Trừ 1 token và cho phép request đi tiếp
	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}
	// Không đủ token -> Từ chối request
	return false
}

// ServeHTTP đóng vai trò làm Middleware kiểm soát luồng HTTP Request
func (rl *RateLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Nếu không đủ token (tần suất vượt quá giới hạn) -> Trả về lỗi 429 Too Many Requests ngay lập tức
	if !rl.Allow() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(429)
		w.Write([]byte(`{"error": "Too Many Requests", "status": 429}`))
		return
	}
	// Router ( DefaultServeMux ) tìm xem handler nào match với URL -> Chuyển tiếp request sang handler đó
	http.DefaultServeMux.ServeHTTP(w, r)
}

// RequestQueue triển khai Hàng đợi đồng bộ sử dụng Buffered Channel
type RequestQueue struct {
	sem   chan struct{} // Cửa sổ giới hạn số lượng worker xử lý đồng thời
	queue chan struct{} // Hàng đợi chứa các request chờ trong RAM
	next  http.Handler
}

func NewRequestQueue(maxWorkers int, maxQueue int, next http.Handler) *RequestQueue {
	return &RequestQueue{
		sem:   make(chan struct{}, maxWorkers),
		queue: make(chan struct{}, maxQueue),
		next:  next,
	}
}

func (rq *RequestQueue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Thử lấy slot worker trực tiếp nếu còn rảnh (không tốn thời gian chờ)
	select {
	case rq.sem <- struct{}{}:
		defer func() { <-rq.sem }()
		rq.next.ServeHTTP(w, r)
		return
	default:
		// Các worker đều đang bận -> Chuyển sang xếp hàng
	}

	// 2. Thử vào hàng đợi (Queue)
	select {
	case rq.queue <- struct{}{}:
		defer func() { <-rq.queue }()
	default:
		// Hàng đợi ĐÃ ĐẦY -> Phản hồi ngay HTTP 503 (Fail-Fast)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"error": "Queue Full - Service Unavailable", "status": 503}`))
		return
	}

	// 3. Đang trong Hàng đợi: Chờ worker nhả slot HOẶC Client hủy kết nối / Timeout
	select {
	case rq.sem <- struct{}{}:
		defer func() { <-rq.sem }()
		rq.next.ServeHTTP(w, r)
	case <-r.Context().Done():
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(499)
		w.Write([]byte(`{"error": "Client Closed Request While Queued", "status": 499}`))
		return
	}
}

// type limitListener struct {
// 	net.Listener
// 	sem chan struct{} // Mỗi phần tử chiếm 0 byte trong mem
// }

// func LimitListener(l net.Listener, n int) net.Listener {
// 	return &limitListener{
// 		Listener: l,
// 		sem: make(chan struct{}, n),
// 	}
// }

// func (l *limitListener) Accept() (net.Conn, error) {
// 	// 1. Chiếm 1 slot kết nối ( Nếu channel đầy -> block)
// 	l.sem<-struct{}{}
// 	// 2. Chấp nhận kết nối từ listener gốc
// 	c, err := l.Listener.Accept()
// 	if err != nil {
// 		<-l.sem // Giải phóng slot nếu bị lỗi
// 		return nil, err
// 	}
// 	// 3. Trả về 1 connection wrapper
// 	return &LimitListenerConn{
// 		Conn: c,
// 		sem: l.sem,
// 	}, nil
// }

// // Định nghĩa Connection Wrapper
// type LimitListenerConn struct {
// 	net.Conn
// 	releaseOnce sync.Once
// 	// Đảm báo chỉ giải phóng slot 1 lần ! cho mỗi connection
// 	sem chan struct{}
// }

// // Override hàm close để giải phóng slot trong channel

// func (c * LimitListenerConn) Close() error {
// 	// Gọi Close của connection gốc
// 	err := c.Conn.Close()

// 	// Giải phóng 1 slot trong sem, nhưng đảm bảo chỉ chạy 1 lần
// 	c.releaseOnce.Do(func() {
// 		<- c.sem
// 	})
// 	return err
// }

func main() {
	algorithm.SetStrategy(&algorithm.RoundRobin{})
	// algorithm.SetStrategy(&algorithm.WeightedRR{})
	// algorithm.SetStrategy(&algorithm.LoadBalancer{})

	http.HandleFunc(
		"/nsmf-pdusession/v1/sm-contexts",
		ForwardToPDU,
	)

	http.HandleFunc(
		"/instances",
		handler.GetInstances,
	)

	http.HandleFunc(
		"/set-weight",
		handler.SetWeight,
	)

	http.HandleFunc(
		"/list-sessions",
		ListSessionsForward,
	)

	log.Println(
		"Gateway started: 8080",
	)

	go func() {
		for {
			health.CheckAllInstances()
			time.Sleep(registry.DefaultInterval)
		}
	}()

	go func() {
		for {
			if algorithm.IsLoadBalancer() {
				health.UpdateAllMetrics()
			}
			time.Sleep(registry.DefaultInterval)
		}
	}()

	go registry.ServiceDiscovery()
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen on :8080: %v", err)
	}
	// LimitListener giới hạn số lượng kết nối đồng thời
	// limitedListener := LimitListener(listener, 20)
	// Rate Limiter & Request Queue (Tạm tắt để bypass điểm nghẽn Mutex/Channel khi test TPS tối đa)
	// limiter := NewRateLimiter(200000, 1000000)
	// reqQueue := NewRequestQueue(100000, 500000, limiter)
	h2s := &http2.Server{
		MaxConcurrentStreams: 10000,
		MaxReadFrameSize:     1048576,
		IdleTimeout:          120 * time.Second,
	}
	server := &http.Server{
		Addr:         ":8080",
		Handler:      h2c.NewHandler(http.DefaultServeMux, h2s),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	log.Println("Gateway started: 8080")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
