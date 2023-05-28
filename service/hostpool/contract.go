package hostpool

import "context"

type (
	UsecaseHostpool interface {
		AddHost(ctx context.Context, host string) (err error)
		RemoveHost(ctx context.Context, host string) (err error)
	}
)
