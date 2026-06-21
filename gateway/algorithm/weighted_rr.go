package algorithm

import (
	"gateway/models"
	"sync"
	"sync/atomic"
)

type WeightedRR struct {
	mutex sync.Mutex
}

func (w *WeightedRR) Select(instances []*models.Instance) *models.Instance {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if len(instances) == 0 {
		return nil
	}
	totalWeight := int32(0)
	for _, ins := range instances {
		totalWeight += atomic.LoadInt32(&ins.Weight)
	}

	for i := range instances {
		instances[i].CurrentWeight += atomic.LoadInt32(&instances[i].Weight)
	}

	maxCurrentWeight := int32(0)
	indexMax := 0
	for i := range instances {
		if instances[i].CurrentWeight > maxCurrentWeight {
			maxCurrentWeight = instances[i].CurrentWeight
			indexMax = i
		}
	}

	instances[indexMax].CurrentWeight -= totalWeight

	return instances[indexMax]
}
