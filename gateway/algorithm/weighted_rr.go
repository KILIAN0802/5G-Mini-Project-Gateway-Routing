package algorithm

import (
	"gateway/models"
	"sync"
)

type WeightedRR struct {
	mutex sync.Mutex
}

func (w *WeightedRR) Select(instances []*models.Instance) models.Instance {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	totalWeight := 0
	for _, ins := range instances {
		totalWeight += ins.Weight
	}

	for i := range instances {
		instances[i].CurrentWeight += instances[i].Weight
	}

	maxCurrentWeight := 0
	indexMax := 0
	for i := range instances {
		if instances[i].CurrentWeight > maxCurrentWeight {
			maxCurrentWeight = instances[i].CurrentWeight
			indexMax = i
		}
	}

	instances[indexMax].CurrentWeight -= totalWeight

	return *instances[indexMax]
}
