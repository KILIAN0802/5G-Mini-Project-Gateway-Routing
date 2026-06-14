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
) models.Instance {
	return CurrentStrategy.Select(instances)
}