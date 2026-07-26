package registry

import (
	"gateway/models"
	"log"
	"net"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

var (
	Instance     []*models.Instance
	RegistryMu   sync.RWMutex
	healthyCache atomic.Value
)

func UpdateHealthyCache() {
	RegistryMu.RLock()
	var healthy []*models.Instance
	for i := range Instance {
		if Instance[i].Healthy.Load() {
			healthy = append(healthy, Instance[i])
		}
	}
	RegistryMu.RUnlock()

	if healthy == nil {
		healthy = []*models.Instance{}
	}
	healthyCache.Store(healthy)
}

func GetHealthyInstance() []*models.Instance {
	val := healthyCache.Load()
	if val == nil {
		return nil
	}
	return val.([]*models.Instance)
}

const DefaultInterval = 10 * time.Second

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
				targetURLStr := "http://" + addr + "/create-session"
				parsedURL, _ := url.Parse(targetURLStr)
				inst := &models.Instance{
					ID:           "Instance:" + addr,
					Address:      addr,
					TargetURL:    parsedURL,
					TargetURLStr: targetURLStr,
					Weight:       1,
				}
				inst.Healthy.Store(true)
				Instance = append(Instance, inst)
				log.Println("New instance added: ", addr)
			}	
		}

		var updated []*models.Instance
		for _, inst := range Instance{
			if newInstances[inst.Address] {
				updated = append(updated, inst)
			} else {
				log.Println("Instance removed: ", inst.Address)
			}
		}

		Instance = updated
		RegistryMu.Unlock()
		UpdateHealthyCache()
		time.Sleep(DefaultInterval)
	}
}
	




