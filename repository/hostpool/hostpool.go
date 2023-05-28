package hostpool

import (
	"context"
)

func (r Repository) GetHostListFromPool(ctx context.Context) (res []string, err error) {
	lock.RLock()
	defer lock.RUnlock()

	return pool, nil
}

func (r Repository) GetCurrentIndex(ctx context.Context) (res int, err error) {
	lock.RLock()
	defer lock.RUnlock()

	return index, nil
}

func (r Repository) GetHostListLen(ctx context.Context) (res int, err error) {
	lock.RLock()
	defer lock.RUnlock()

	return len(pool), nil
}

func (r Repository) AppendHost(ctx context.Context, host string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	pool = append(pool, host)

	return nil
}

func (r Repository) RequeueFirstHostToLast(ctx context.Context) (err error) {
	lock.Lock()
	defer lock.Unlock()

	if len(pool) > 0 {
		pool = append(pool[0:], pool[0])
	}

	return nil
}

func (r Repository) RemoveHostByHostAddress(ctx context.Context, host string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	for i, value := range pool { //TODO: can optimize? but you dont usually have a large set of host so...
		if value == host {
			// Remove the element from the roundRobinQueue
			pool = append(pool[:i], pool[i+1:]...)
			index = index % len(pool)
			break
		}
	}

	return nil

}

func (r Repository) IncrementIndex(ctx context.Context, increment int) (res int, err error) {
	lock.Lock()
	defer lock.Unlock()

	index = (index + increment) % len(pool)

	return index, nil
}

func (r Repository) SetIndex(ctx context.Context, newIndex int) (res int, err error) {
	lock.Lock()
	defer lock.Unlock()

	index = newIndex % len(pool)

	return index, nil
}
