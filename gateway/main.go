package main

import (
	"bytes" // Cung cấp các hàm để thao tác với kiểu dữ liệu byte -> Thường dùng để tạo buffer cho việc đọc / ghi dữ liệu
	"io"    // Cung cấp các interface chuẩn để đọc/ghi dữ liệu
	"context"
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	// Cho phép:
	// Tạo web server ( http.ListenAndServe)
	// Gửi request HTTP (http.Get, http.Post)
	// Xử lý request/response qua http.Handler
	"gateway/algorithm"
	"gateway/handler"
	"gateway/health"
	"gateway/registry"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Notice
var pduClient = &http.Client{
	Timeout: 90 * time.Second,
	Transport: &http2.Transport{
		AllowHTTP: true,
		DialTLSContext: func(ctx context.Context, network, addr string, cfg *tls.Config) (net.Conn, error){
			var d net.Dialer
			return d.DialContext(ctx, network, addr)
		},
	},
}

func ForwardToPDU(
	w http.ResponseWriter,
	r *http.Request,
) {
	bodyBytes, err :=
		io.ReadAll(r.Body)

	if err != nil {
		http.Error(
			w,
			"Can not read body",
			500,
		)
		return
	}
	selected := algorithm.SelectBackend(registry.GetHealthyInstance())
	if selected == nil {
		http.Error(
			w,
			"NO_BACKEND_AVAILABLE",
			503,
		)
		return
	}

	isLB := algorithm.IsLoadBalancer()
	if isLB {
		atomic.AddInt32(&selected.ActiveRequest, 1)
	}
	defer func() {
		if isLB {
			atomic.AddInt32(&selected.ActiveRequest, -1)
		}
	}()

	log.Printf(
		"Gateway route to %s",
		selected.ID,
	)

	req, err :=
		http.NewRequest(
			"POST",
			"http://"+selected.Address+"/create-session",
			bytes.NewBuffer(bodyBytes),
		)
	if err != nil {
		http.Error(
			w,
			"Error creating request",
			500,
		)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := pduClient.Do(req)
	if err != nil {
		http.Error(
			w,
			"Backend Error",
			500,
		)
		return
	}

	defer resp.Body.Close()

	if _, err := io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Error forwarding response", 500)
		log.Printf(
			"Error forwarding response: %v",
			err,
		)
		return
	}
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

type RateLimiter struct {
	rate       float64
	capacity   float64
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewRateLimiter(rate float64, capacity float64) *RateLimiter {
	return &RateLimiter{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill).Seconds()
	rl.lastRefill = now

	rl.tokens += elapsed * rl.rate
	if rl.tokens > rl.capacity {
		rl.tokens = rl.capacity
	}

	if rl.tokens >= 1.0 {
		rl.tokens -= 1.0
		return true
	}
	return false
}

func (rl *RateLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !rl.Allow() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(429)
		w.Write([]byte(`{"error": "Too Many Requests", "status": 429}`))
		return
	}
	// Router ( DefaultServeMux ) tìm xem handle nào match với URL -> Serve 
	http.DefaultServeMux.ServeHTTP(w, r)
}

type limitListener struct {
	net.Listener
	sem chan struct{} // Mỗi phần tử chiếm 0 byte trong mem
}

func LimitListener(l net.Listener, n int) net.Listener {
	return &limitListener{
		Listener: l,
		sem: make(chan struct{}, n),
	}
}

func (l *limitListener) Accept() (net.Conn, error) {
	// 1. Chiếm 1 slot kết nối ( Nếu channel đầy -> block)
	l.sem<-struct{}{}
	// 2. Chấp nhận kết nối từ listener gốc
	c, err := l.Listener.Accept()
	if err != nil {
		<-l.sem // Giải phóng slot nếu bị lỗi
		return nil, err
	}
	// 3. Trả về 1 connection wrapper 
	return &LimitListenerConn{
		Conn: c, 
		sem: l.sem,
	}, nil
}

// Định nghĩa Connection Wrapper
type LimitListenerConn struct {
	net.Conn
	releaseOnce sync.Once 
	// Đảm báo chỉ giải phóng slot 1 lần ! cho mỗi connection
	sem chan struct{}
}

// Override hàm close để giải phóng slot trong channel

func (c * LimitListenerConn) Close() error {
	// Gọi Close của connection gốc
	err := c.Conn.Close()

	// Giải phóng 1 slot trong sem, nhưng đảm bảo chỉ chạy 1 lần
	c.releaseOnce.Do(func() {
		<- c.sem
	})
	return err
}


func main() {
	// algorithm.SetStrategy(&algorithm.RoundRobin{})
	// algorithm.SetStrategy(&algorithm.WeightedRR{})
	algorithm.SetStrategy(&algorithm.LoadBalancer{})

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
	limitedListener := LimitListener(listener, 10000)
	// Rate Limiter
	limiter := NewRateLimiter(20000, 10000)
	h2s := &http2.Server{}
	http.Serve(limitedListener, h2c.NewHandler(limiter, h2s))
}
