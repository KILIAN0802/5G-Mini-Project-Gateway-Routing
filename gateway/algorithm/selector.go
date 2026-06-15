package algorithm

import "gateway/models"

var CurrentStrategy Strategy

func SetStrategy(
	strategy Strategy,
) {
	CurrentStrategy = strategy
}

func SelectBackend(
	instances []*models.Instance,
) *models.Instance {
	return CurrentStrategy.Select(instances)
}

func IsLoadBalancer() bool {
	if CurrentStrategy == nil {
		return false
	}
	_, ok := CurrentStrategy.(*LoadBalancer)
	if !ok {
		_, ok = CurrentStrategy.(LoadBalancer)
	}
	return ok
}