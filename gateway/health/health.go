package health

import (
	"gateway/models"
	"gateway/registry"
	"log"
	"net/http"
	"time"
)

func CheckInstance(
	instance *models.Instance,
) {
	client := http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err :=
		client.Get(
			"http://" + instance.Address + "/health",
		)

	if err != nil {
		instance.Healthy = false
		log.Printf(
			"%s healthy=%v",
			instance.ID,
			instance.Healthy,
		)
		return
	}

	resp.Body.Close()

	instance.Healthy =
		resp.StatusCode ==
			200

	log.Printf(
		"%s healthy=%v",
		instance.ID,
		instance.Healthy,
	)
}

func CheckAllInstances() {
	for i := range registry.Instances {
		CheckInstance(
			&registry.Instances[i],
		)
	}
}
