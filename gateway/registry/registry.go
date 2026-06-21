package registry

import (
	"gateway/models"
	"log"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func GetHealthyInstance() []*models.Instance {

	RegistryMu.RLock()
	defer RegistryMu.RUnlock()

	var healthy []*models.Instance // slice sẽ chứa những pdu instance đang hoạt động
	for i := range Instance {
		if Instance[i].Healthy.Load() {
			healthy =
				append(
					healthy,
					Instance[i],
				)// thêm Instance[i] vào health
		}
	}

	return healthy
}

const DefaultInterval = 10 * time.Second

var (
	Instance []*models.Instance
	RegistryMu sync.RWMutex
)

func ServiceDiscovery() {
	for {
		ips, err := net.LookupIP("pdu-session")
		if err != nil{
			log.Println("Lookup error : ", err)
			time.Sleep(DefaultInterval)
			continue
		}

		newInstances := make(map[string] bool)
		for _, ip := range ips {
			addr := ip.String() + ":9001"
			newInstances[addr] = true
		}

		RegistryMu.Lock()
		// Thêm IP mới

		for addr := range newInstances {
			found := false
			for _, inst := range Instance {
				if inst.Address == addr {
					found = true
					break
				}
			}

			if !found {
				Instance = append(Instance, &models.Instance{
					ID: "Instance:"+ addr,
					Address: addr,
					Healthy: atomic.Bool{},
					Weight: 1,
				})
				log.Println("New instance added: ", addr)
			}	
			}

			var updated []*models.Instance
			for _, inst := range Instance{
				if newInstances[inst.Address] {
				updated = append(updated, inst)
			}else{
				log.Println("Instance removed: ", inst.Address)
			}
		}

		Instance = updated
		RegistryMu.Unlock()
		time.Sleep(DefaultInterval)
	}
}
	




