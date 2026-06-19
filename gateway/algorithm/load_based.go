package algorithm

import (
	"gateway/models"
	"sync/atomic"
)

type LoadBalancer struct{}

func (l LoadBalancer) Select(
	instances []*models.Instance,
) *models.Instance {
	if len(instances) == 0 {
		return nil
	}

	selected := instances[0]

	for _, ins := range instances {
		if atomic.LoadInt32(&ins.ActiveRequest) < atomic.LoadInt32(&selected.ActiveRequest) {
			selected = ins
		}
	}
	return selected
}