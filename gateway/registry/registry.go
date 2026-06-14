package registry

import (
	"gateway/models"
)

func GetHealthyInstance() []*models.Instance {
	var healthy []*models.Instance

	for i := range Instances {
		if Instances[i].Healthy {
			healthy =
				append(
					healthy,
					&Instances[i],
				)
		}
	}

	return healthy
}

var Instances = []models.Instance{
	{
		ID:      "pdu-1",
		Address: "localhost:9001",
		Weight: 3,
	},
	{
		ID:      "pdu-2",
		Address: "localhost:9002",
		Weight: 1,
	},
	{
		ID:      "pdu-3",
		Address: "localhost:9003",
		Weight: 6,
	},
	{
		ID:      "pdu-4",
		Address: "localhost:9004",
		Weight: 4,
	},
}
