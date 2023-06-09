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

func (r Repository) GetHostListLength(ctx context.Context) (res int, err error) {
	lock.RLock()
	defer lock.RUnlock()

	return len(pool), nil
}

func (r Repository) AppendHost(ctx context.Context, host string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	// ensure no duplicate host, still need to deal with http:// or https:// considered different tho...
	if poolMap[host] {
		return nil
	}

	pool = append(pool, host)
	poolMap[host] = true

	return nil
}

func (r Repository) RequeueFirstHostToLast(ctx context.Context) (err error) {
	lock.Lock()
	defer lock.Unlock()

	if len(pool) > 1 {
		pool = append(pool[1:], pool[0])
	}

	return nil
}

func (r Repository) RemoveHostByHostAddress(ctx context.Context, host string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	for i, value := range pool { //TODO: can optimize? but you dont usually have a large set of hostpool so...
		if value == host {
			// Remove the element from the roundRobinQueue
			pool = append(pool[:i], pool[i+1:]...)
			poolMap[host] = false
			if len(pool) > 0 {
				index = index % len(pool)
			} else {
				index = 0
			}
			break
		}
	}

	return nil

}

func (r Repository) IncrementIndex(ctx context.Context, increment int) (res int, err error) {
	lock.Lock()
	defer lock.Unlock()

	if len(pool) > 0 {
		index = (index + increment) % len(pool)
	} else {
		index = 0
	}

	return index, nil
}

func (r Repository) SetIndex(ctx context.Context, newIndex int) (res int, err error) {
	lock.Lock()
	defer lock.Unlock()

	if len(pool) > 0 {
		index = newIndex % len(pool)
	} else {
		index = 0
	}

	return index, nil
}
