package hostpool

import "context"

type (
	RepoHostpool interface {
		GetHostListFromPool(ctx context.Context) (res []string, err error)
		GetCurrentIndex(ctx context.Context) (index int, err error)
		GetHostListLen(ctx context.Context) (res int, err error)
		AppendHost(ctx context.Context, host string) (err error)
		RequeueFirstHostToLast(ctx context.Context) (err error)
		RemoveHostByHostAddress(ctx context.Context, host string) (err error)
		IncrementIndex(ctx context.Context, increment int) (res int, err error)
		SetIndex(ctx context.Context, newIndex int) (res int, err error)
	}
)
