package forwarder

import "context"

type (
	UsecaseHostpool interface {
		GetHost(ctx context.Context) (res string, err error)
	}
)
