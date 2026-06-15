package health

import (
	"encoding/json"
	"log"
	"net/http"
	"gateway/models"
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
		instance.ActiveRequest,
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

	instance.ActiveRequest = metrics.ActiveRequests
}

func UpdateAllMetrics(instances []models.Instance){
	for i := range instances{
		UpdateMetrics(&instances[i])
	}
}