package cronjob

import (
	"github.com/loadbalancer/pkg/config"
	"github.com/loadbalancer/service/hostpool"
	"gopkg.in/robfig/cron.v2"
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
		hostpoolService, err := hostpool.GetInstance(c)
		if err != nil {
			errFinal = err
			return
		}

		m = Module{
			config:     c,
			cronModule: cron.New(),
			service: service{
				hostpool: hostpoolService,
			},
		}
	})

	return m, errFinal
}
