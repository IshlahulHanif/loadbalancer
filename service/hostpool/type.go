package hostpool

import "github.com/loadbalancer/pkg/config"

type (
	Service struct {
		config  config.Config
		usecase usecase
	}

	usecase struct {
		hostpool   UsecaseHostpool
		poolclient UsecasePoolClient
	}
)
