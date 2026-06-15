package models

type MetricsResponse struct {
	InstanceID string `json:"instanceID"`
	ActiveRequests int `json:"activeRequests"`
}