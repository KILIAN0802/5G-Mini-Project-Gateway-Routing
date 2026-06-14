package algorithm

import "gateway/models"

var counter int

type RoundRobin struct{}

func (r RoundRobin) Select(healthy []*models.Instance) models.Instance {
	if len(healthy) == 0 {
		return models.Instance{}
	}
	index := counter % len(healthy)

	counter++
	return *healthy[index]
}
