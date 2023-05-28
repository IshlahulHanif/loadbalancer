package hostpool

import (
	"github.com/IshlahulHanif/poneglyph"
	"github.com/loadbalancer/pkg/config"
	"github.com/loadbalancer/repository/hostpool"
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
		hostpoolRepo, err := hostpool.GetInstance(c)
		if err != nil {
			errFinal = poneglyph.Trace(err)
			return
		}

		m = Usecase{
			repo: repository{
				hostpool: hostpoolRepo,
			},
		}
	})

	return m, errFinal
}
