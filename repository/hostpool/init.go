package hostpool

import (
	"github.com/loadbalancer/pkg/config"
	"sync"
)

var (
	m    Repository
	once sync.Once

	// hostpool pool data related
	lock    sync.RWMutex
	pool    []string
	poolMap map[string]bool
	index   int
)

func GetInstance(c config.Config) (Repository, error) {
	var (
		errFinal error
	)

	once.Do(func() {
		// append all hostpool from config
		for _, host := range c.HostList {
			pool = append(pool, host)
		}

		poolMap = make(map[string]bool)

		m = Repository{}
	})

	return m, errFinal
}
