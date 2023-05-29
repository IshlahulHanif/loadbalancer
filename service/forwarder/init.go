package forwarder

import (
	"github.com/IshlahulHanif/poneglyph"
	"github.com/loadbalancer/pkg/config"
	"github.com/loadbalancer/usecase/hostpool"
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

		m = Service{
			config: c,
			usecase: usecase{
				hostpool: hostpoolUsecase,
			},
		}
	})

	return m, errFinal
}
