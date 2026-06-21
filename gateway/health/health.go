package health

import (
	"gateway/models"
	"gateway/registry"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

var healthClient = &http.Client{
	Timeout: 2 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     30 * time.Second,
	},
}

func CheckInstance(
	instance *models.Instance,
) {
	resp, err :=
		healthClient.Get(
			"http://" + instance.Address + "/health",
		)



	healthy := false
	if err == nil {
		healthy = (resp.StatusCode == 200)
	}
	if resp != nil {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}

	instance.Healthy.Store(healthy)
	log.Printf("%s healthy=%v", instance.ID, healthy)
}


func CheckAllInstances() {
	registry.RegistryMu.RLock()
	instances := make([]*models.Instance, len(registry.Instance))
	copy(instances, registry.Instance)
	registry.RegistryMu.RUnlock()

	var wg sync.WaitGroup
	for _, inst := range instances {
		wg.Add(1)
		go func(i *models.Instance) {
			defer wg.Done()
			CheckInstance(i)
		}(inst)
	}
	wg.Wait()
}
