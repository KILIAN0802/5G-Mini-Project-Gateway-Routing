package health

import (
	"encoding/json"
	"gateway/models"
	"gateway/registry"
	"log"
	"net/http"
	// net: định nghĩa các kiểu dữ liệu về internet như địa chỉ ip, ...
	// net/http: định nghĩa các phương thức http như GET, POST, PUT, DELETE
	"sync/atomic"
)

func UpdateMetrics(
	instance *models.Instance,
) {
	resp, err := http.Get(
		"http://"+instance.Address+"/metrics",
	)
	log.Printf(
		"%s active=%d",
		instance.ID,
		atomic.LoadInt32(&instance.ActiveRequest),
	)

	if err !=nil{
		return
	}

	var metrics models.MetricsResponse

	err = json.NewDecoder(
		resp.Body,
	).Decode(
		&metrics,
	)

	if err !=nil{
		return
	}

	atomic.StoreInt32(&instance.ActiveRequest, int32(metrics.ActiveRequests))
}

func UpdateAllMetrics(){
	registry.RegistryMu.RLock() // Chỉ cho phép đọc không cho ghi
	defer registry.RegistryMu.RUnlock()// Bảo vệ registry.Instance
	for i := range registry.Instance{
		UpdateMetrics(registry.Instance[i])
	}
}