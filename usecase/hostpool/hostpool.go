package hostpool

import (
	"context"
	"errors"
	"github.com/IshlahulHanif/poneglyph"
)

func (u Usecase) GetHost(ctx context.Context) (res string, err error) {
	pool, err := u.repo.hostpool.GetHostListFromPool(ctx)
	if err != nil {
		return "", poneglyph.Trace(err)
	}

	length := len(pool)
	if length == 0 {
		// something wrong
		return "", poneglyph.Trace(errors.New("hostpool empty"))
	}

	currentIdx, err := u.repo.hostpool.GetCurrentIndex(ctx)
	if err != nil {
		return "", poneglyph.Trace(err)
	}

	if length-1 < currentIdx {
		// index overload, restart to 0
		currentIdx, err = u.repo.hostpool.SetIndex(ctx, 0)
		if err != nil {
			return "", poneglyph.Trace(err)
		}

		// re-arrange hostpool list to avoid overwork on first hostpool
		err = u.repo.hostpool.RequeueFirstHostToLast(ctx)
		if err != nil {
			return "", poneglyph.Trace(err)
		}

		return "", poneglyph.Trace(err)
	}

	res = pool[currentIdx]
	currentIdx, err = u.repo.hostpool.IncrementIndex(ctx, 1)
	if err != nil {
		return "", poneglyph.Trace(err)
	}

	return res, nil
}

func (u Usecase) AddHost(ctx context.Context, host string) (err error) {
	err = u.repo.hostpool.AppendHost(ctx, host)
	if err != nil {
		return poneglyph.Trace(err)
	}

	return nil
}

func (u Usecase) RemoveHost(ctx context.Context, host string) (err error) {
	err = u.repo.hostpool.RemoveHostByHostAddress(ctx, host)
	if err != nil {
		return poneglyph.Trace(err)
	}

	return nil
}
