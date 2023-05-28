package poolclient

import (
	"github.com/loadbalancer/pkg/config"
	"sync"
)

var (
	m    Usecase
	once sync.Once
)

func GetInstance(c config.Config) (Usecase, error) {
	var (
		errFinal error
	)

	once.Do(func() {
		m = Usecase{
			config: c,
		}
	})

	return m, errFinal
}
