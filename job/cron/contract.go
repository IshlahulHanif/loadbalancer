package cronjob

import (
	"context"
	hostEntity "github.com/loadbalancer/entity/host"
)

type (
	hostpoolMethod interface {
		HealthCheckAllHost(ctx context.Context) (resp hostEntity.HealthCheckAllHostResp, err error)
	}
)
