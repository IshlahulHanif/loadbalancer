package hostpool

import (
	"github.com/IshlahulHanif/poneglyph"
	"github.com/loadbalancer/pkg/config"
	"github.com/loadbalancer/usecase/hostpool"
	"github.com/loadbalancer/usecase/poolclient"
	"sync"
)

var (
	m    Service
	once sync.Once
)

func GetInstance(c config.Config) (Service, error) {
	var (
		errFinal error
	)

	once.Do(func() {
		hostpoolUsecase, err := hostpool.GetInstance(c)
		if err != nil {
			errFinal = poneglyph.Trace(err)
			return
		}

		poolclientUsecase, err := poolclient.GetInstance(c)
		if err != nil {
			errFinal = poneglyph.Trace(err)
			return
		}

		m = Service{
			config: c,
			usecase: usecase{
				hostpool:   hostpoolUsecase,
				poolclient: poolclientUsecase,
			},
		}
	})

	return m, errFinal
}
