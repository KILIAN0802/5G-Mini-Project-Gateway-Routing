package main

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"gateway/algorithm"
	"gateway/handler"
	"gateway/health"
	"gateway/registry"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 32768)
		return &b
	},
}

const clientPoolSize = 32

var (
	pduClientPool []*http.Client
	pduClientIdx  uint64
)

func initPDUClientPool() {
	pduClientPool = make([]*http.Client, clientPoolSize)
	for i := 0; i < clientPoolSize; i++ {
		pduClientPool[i] = &http.Client{
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
						_ = tcpConn.SetReadBuffer(1024 * 1024)
						_ = tcpConn.SetWriteBuffer(1024 * 1024)
					}
					return conn, nil
				},
				StrictMaxConcurrentStreams: false,
				ReadIdleTimeout:            30 * time.Second,
				PingTimeout:                15 * time.Second,
			},
		}
	}
}

func getPDUClient() *http.Client {
	idx := atomic.AddUint64(&pduClientIdx, 1)
	return pduClientPool[idx%clientPoolSize]
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

	var req *http.Request
	if selected.TargetURL != nil {
		req = &http.Request{
			Method:        "POST",
			URL:           selected.TargetURL,
			Header:        r.Header,
			Body:          r.Body,
			Host:          selected.Address,
			ContentLength: r.ContentLength,
		}
	} else {
		var err error
		req, err = http.NewRequestWithContext(r.Context(), "POST", "http://"+selected.Address+"/create-session", r.Body)
		if err != nil {
			http.Error(w, "Error creating request", 500)
			return
		}
		req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	}

	resp, err := getPDUClient().Do(req)
	if err != nil {
		log.Printf("ForwardToPDU error forwarding to %s: %v", selected.Address, err)
		http.Error(w, "Backend Error: "+err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	bufPtr := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(bufPtr)
	io.CopyBuffer(w, resp.Body, *bufPtr)
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

	resp, err := getPDUClient().Do(req)
	if err != nil {
		log.Printf("ListSessionsForward error forwarding to %s: %v", selected.Address, err)
		http.Error(w, "Backend Error: "+err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	bufPtr := bufferPool.Get().(*[]byte)
	defer bufferPool.Put(bufPtr)
	io.CopyBuffer(w, resp.Body, *bufPtr)
}

func main() {
	initPDUClientPool()
	registry.UpdateHealthyCache()
	algorithm.SetStrategy(&algorithm.RoundRobin{})

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

	h2s := &http2.Server{
		MaxConcurrentStreams:         50000,
		MaxReadFrameSize:             1048576,
		IdleTimeout:                  120 * time.Second,
		MaxUploadBufferPerStream:     65535 * 32,
		MaxUploadBufferPerConnection: 65535 * 64,
	}
	server := &http.Server{
		Addr:         ":8080",
		Handler:      h2c.NewHandler(http.DefaultServeMux, h2s),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	log.Println("Gateway started on :8080")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

