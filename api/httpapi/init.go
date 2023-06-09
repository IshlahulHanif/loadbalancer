package httpapi

import (
	"github.com/loadbalancer/pkg/config"
	"github.com/loadbalancer/service/forwarder"
	"github.com/loadbalancer/service/hostpool"
	"sync"
)

var (
	m    Module
	once sync.Once
)

func GetInstance(c config.Config) (Module, error) {
	var (
		errFinal error
	)

	once.Do(func() {
		forwarderService, err := forwarder.GetInstance(c)
		if err != nil {
			errFinal = err
			return
		}

		hostpoolService, err := hostpool.GetInstance(c)
		if err != nil {
			errFinal = err
			return
		}

		m = Module{
			service: service{
				forwarder: forwarderService,
				hostpool:  hostpoolService,
			},
		}
	})

	return m, errFinal
}
