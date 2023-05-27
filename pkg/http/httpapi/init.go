package httpapi

import (
	"github.com/loadbalancer/pkg/config"
	"github.com/loadbalancer/usecase/loadbalancer"
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
		loadbalancerUsecase, err := loadbalancer.GetInstance(c)
		if err != nil {
			errFinal = err
			return
		}

		m = Module{
			usecase: usecase{
				loadbalancer: loadbalancerUsecase,
			},
		}
	})

	return m, errFinal
}
