package httpapi

import (
	"context"
	"github.com/loadbalancer/entity/forwarder"
	"github.com/loadbalancer/entity/host"
)

type (
	forwarderMethod interface {
		ForwardRequest(ctx context.Context, req forwarder.ForwardRequestReq) (resp forwarder.ForwardRequestResp, err error)
	}

	hostpoolMethod interface {
		ManageHost(ctx context.Context, req host.ManageHostReq) (err error)
	}
)
