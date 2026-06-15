package main

import (
	"bytes" // Cung cấp các hàm để thao tác với kiểu dữ liệu byte -> Thường dùng để tạo buffer cho việc đọc / ghi dữ liệu
	"io"    // Cung cấp các interface chuẩn để đọc/ghi dữ liệu
	"log"
	"net/http"

	// Cho phép:
	// Tạo web server ( http.ListenAndServe)
	// Gửi request HTTP (http.Get, http.Post)
	// Xử lý request/response qua http.Handler
	"gateway/algorithm"
	"gateway/handler"
	"gateway/health"
	"gateway/models"
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
	if selected == (models.Instance{}) {
		http.Error(
			w,
			"NO_BACKEND_AVAILABLE",
			500,
		)
		return
	}

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
	algorithm.SetStrategy(&algorithm.RoundRobin{})
	// algorithm.SetStrategy(&algorithm.WeightedRR{})

	http.HandleFunc(
		"/nsmf-pdusession/v1/sm-contexts",
		ForwardToPDU,
	)

	http.HandleFunc(
		"/instances",
		handler.GetInstances,
	)

	log.Println(
		"Gateway started: 8080",
	)

	go func() {
		for {
			health.CheckAllInstances()
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			health.UpdateAllMetrics(
				registry.Instances,
			)
			time.Sleep(1 * time.Second)
		}
	}()

	http.ListenAndServe(
		":8080",
		nil,
	)
}
