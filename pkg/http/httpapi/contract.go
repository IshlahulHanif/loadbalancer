package httpapi

import "context"

type (
	loadbalancerMethod interface {
		BouncerApi(ctx context.Context, req interface{}) (res interface{}, err error)
	}
)
