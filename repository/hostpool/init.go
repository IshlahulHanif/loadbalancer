package hostpool

import (
	"github.com/loadbalancer/pkg/config"
	"sync"
)

var (
	m    Repository
	once sync.Once

	// host pool data related
	lock  sync.RWMutex
	pool  []string
	index int
)

func GetInstance(c config.Config) (Repository, error) {
	var (
		errFinal error
	)

	once.Do(func() {
		// append all host from config
		for _, host := range c.HostList {
			pool = append(pool, host)
		}

		m = Repository{}
	})

	return m, errFinal
}
