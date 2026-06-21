package algorithm

import (
	"gateway/models"
	"sync/atomic"
)

var counter uint64

type RoundRobin struct{}

func (r RoundRobin) Select(healthy []*models.Instance) *models.Instance {
	if len(healthy) == 0 {
		return nil
	}
	val := atomic.AddUint64(&counter, 1) - 1
	index := val % uint64(len(healthy))
	return healthy[index]
}
