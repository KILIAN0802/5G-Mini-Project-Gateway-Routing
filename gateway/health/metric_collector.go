package health

import (
	"encoding/json"
	"gateway/models"
	"gateway/registry"
	"log"
	"net/http"
	// net: định nghĩa các kiểu dữ liệu về internet như địa chỉ ip, ...
	// net/http: định nghĩa các phương thức http như GET, POST, PUT, DELETE
	"sync"
	"sync/atomic"
	"time"
)

var metricClient = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	},
}

func UpdateMetrics(
	instance *models.Instance,
) {
	resp, err := metricClient.Get(
		"http://"+instance.Address+"/metrics",
	)
	log.Printf(
		"%s active=%d",
		instance.ID,
		atomic.LoadInt32(&instance.ActiveRequest),
	)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	var metrics models.MetricsResponse

	err = json.NewDecoder(
		resp.Body,
	).Decode(
		&metrics,
	)
	if err != nil {
		return
	}

	atomic.StoreInt32(&instance.ActiveRequest, int32(metrics.ActiveRequests))
}

func UpdateAllMetrics() {
	registry.RegistryMu.RLock()
	instances := make([]*models.Instance, len(registry.Instance))
	copy(instances, registry.Instance)
	registry.RegistryMu.RUnlock()

	var wg sync.WaitGroup
	for _, inst := range instances {
		wg.Add(1)
		go func(i *models.Instance) {
			defer wg.Done()
			UpdateMetrics(i)
		}(inst)
	}
	wg.Wait()
}