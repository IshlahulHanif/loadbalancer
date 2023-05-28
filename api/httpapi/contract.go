package httpapi

import (
	"context"
	"github.com/loadbalancer/entity/forwarder"
)

type (
	forwarderMethod interface {
		ForwardRequest(ctx context.Context, req forwarder.ForwardRequestReq) (resp forwarder.ForwardRequestResp, err error)
	}
)
