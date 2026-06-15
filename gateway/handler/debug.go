package handler

import (
	"encoding/json"
	"gateway/registry"
	"net/http"
)

type InstanceDebugResponse struct {
	ID             string `json:"id"`
	ActiveRequests int    `json:"activeRequests"`
}

func GetInstances(w http.ResponseWriter, r *http.Request) {
	instances := registry.Instances
	resp := make([]InstanceDebugResponse, len(instances))
	for i, inst := range instances {
		resp[i] = InstanceDebugResponse{
			ID:             inst.ID,
			ActiveRequests: inst.ActiveRequest,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}