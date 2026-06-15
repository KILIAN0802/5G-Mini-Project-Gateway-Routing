package algorithm

import "gateway/models"

type LoadBalancer struct{}

func (l LoadBalancer) Select(
	instances []*models.Instance,
) *models.Instance {
	if len(instances) == 0 {
		return nil
	}

	selected := instances[0]

	for _, ins := range instances {
		if ins.ActiveRequest < selected.ActiveRequest {
			selected = ins
		}
	}
	return selected
}