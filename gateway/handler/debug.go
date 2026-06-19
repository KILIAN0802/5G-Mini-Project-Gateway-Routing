package handler

import (
	"encoding/json"
	"gateway/registry"
	"gateway/algorithm"
	"net/http"
	"sync/atomic"
	"strconv"
)

type InstanceDebugResponse struct {
	ID             string `json:"id"`
	ActiveRequests int    `json:"activeRequests"`
	Weight         int    `json:"weight"`
}

func GetInstances(w http.ResponseWriter, r *http.Request) {
	registry.RegistryMu.RLock()
	defer registry.RegistryMu.RUnlock()
	instances := registry.Instance
	resp := make([]InstanceDebugResponse, len(instances))
	for i, inst := range instances {
		resp[i] = InstanceDebugResponse{
			ID:             inst.ID,
			ActiveRequests: int(atomic.LoadInt32(&inst.ActiveRequest)),
			Weight:         inst.Weight,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func SetWeight(w http.ResponseWriter, r *http.Request) {
	// Kiểm tra xem thuật toán hiện tại có phải là Weighted Round Robin hay không
	if !algorithm.IsWeightedRR() {
		http.Error(w, "Weight configuration is only allowed in Weighted Round Robin strategy", http.StatusForbidden)
		return
	}

	addr := r.URL.Query().Get("address")
	weightStr := r.URL.Query().Get("weight")
	
	weight, err := strconv.Atoi(weightStr)
	if err != nil || weight <= 0 {
		http.Error(w, "Invalid weight value", http.StatusBadRequest)
		return
	}

	registry.RegistryMu.Lock()
	defer registry.RegistryMu.Unlock()

	found := false
	for _, inst := range registry.Instance {
		if inst.Address == addr {
			inst.Weight = weight
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Instance not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("Weight updated successfully"))
}