package forwarder

import "context"

type (
	UsecaseHostpool interface {
		GetHostDequeueRoundRobin(ctx context.Context) (res string, err error)
		RemoveHost(ctx context.Context, host string) (err error)
	}
)
