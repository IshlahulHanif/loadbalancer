package hostpool

import (
	"context"
	"errors"
	"github.com/IshlahulHanif/poneglyph"
	"github.com/loadbalancer/entity/host"
)

func (s Service) ManageHost(ctx context.Context, req host.ManageHostReq) (err error) {
	switch req.Operation {
	case host.Add:
		for _, hostIP := range req.Data {
			err = s.usecase.hostpool.AddHost(ctx, hostIP)
			if err != nil {
				return poneglyph.Trace(err)
			}
		}
	case host.Remove:
		for _, hostIP := range req.Data {
			err = s.usecase.hostpool.RemoveHost(ctx, hostIP)
			if err != nil {
				return poneglyph.Trace(err)
			}
		}
	case host.HealthCheck:
		//TODO: implement
	default:
		return poneglyph.Trace(errors.New("unknown host operation"))
	}

	return nil
}
