package hostpool

import "context"

type (
	UsecaseHostpool interface {
		AddHost(ctx context.Context, host string) (err error)
		RemoveHost(ctx context.Context, host string) (err error)
		GetHostListFromPool(ctx context.Context) (res []string, err error)
	}
	UsecasePoolClient interface {
		PingHost(ctx context.Context, host string) (isPingSuccess bool)
	}
)
