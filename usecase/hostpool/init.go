package hostpool

import (
	"github.com/loadbalancer/pkg/config"
	"sync"
)

var (
	m    Usecase
	once sync.Once
	lock sync.Mutex
)

func GetInstance(c config.Config) (Usecase, error) {
	var (
		errFinal error
	)

	once.Do(func() {
		// append all host from config
		var hostList = make([]string, 0)
		for _, host := range c.HostList {
			hostList = append(hostList, host)
		}

		queue := RoundRobinQueue{
			queue: hostList,
			index: 0,
		}

		m = Usecase{
			roundRobinQueue: queue, //TODO: move the sync lock to outside
		}
	})

	return m, errFinal
}
