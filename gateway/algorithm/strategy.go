package algorithm

import (
	"gateway/models"
)

type Strategy interface {
	Select(instances []*models.Instance) models.Instance
}