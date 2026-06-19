package main

import (
	"bytes" // Cung cấp các hàm để thao tác với kiểu dữ liệu byte -> Thường dùng để tạo buffer cho việc đọc / ghi dữ liệu
	"io"    // Cung cấp các interface chuẩn để đọc/ghi dữ liệu
	"log"
	"net/http"
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
)

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
			500,
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

	resp, err :=
		http.Post(
			"http://"+selected.Address+"/create-session",
			"application/json",
			bytes.NewBuffer(bodyBytes),
		)

	if err != nil {
		http.Error(
			w,
			"Backend Error",
			500,
		)
		return
	}

	defer resp.Body.Close()

	responseBody, _ := io.ReadAll(resp.Body)

	w.Write(responseBody)
}

func main() {
	// algorithm.SetStrategy(&algorithm.RoundRobin{})
	algorithm.SetStrategy(&algorithm.WeightedRR{})
	//algorithm.SetStrategy(&algorithm.LoadBalancer{})

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
	http.ListenAndServe(
		":8080",
		nil,
	)
}
