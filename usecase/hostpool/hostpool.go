package hostpool

import (
	"context"
	"fmt"
	"github.com/IshlahulHanif/poneglyph"
)

func (u Usecase) GetHost(ctx context.Context) (res string, err error) {
	host, err := u.Dequeue(ctx)
	if err != nil {
		return "", poneglyph.Trace(err)
	}

	u.Enqueue(ctx, host)
	return host, nil
}

func (u Usecase) AddHost(ctx context.Context, host string) (err error) {
	u.Enqueue(ctx, host)
	return nil
}

func (u Usecase) RemoveHost(ctx context.Context, host string) (err error) {
	u.Remove(ctx, host)
	return nil
}

// TODO: below funtions should be on repo
func (u Usecase) Enqueue(ctx context.Context, element string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	u.roundRobinQueue.queue = append(u.roundRobinQueue.queue, element)

	return nil
}

func (u Usecase) Dequeue(ctx context.Context) (res string, err error) {
	lock.Lock()
	defer lock.Unlock()

	length := len(u.roundRobinQueue.queue)
	if length == 0 {
		err = fmt.Errorf("roundRobinQueue is empty")
		return "", poneglyph.Trace(err)
	}

	element := u.roundRobinQueue.queue[u.roundRobinQueue.index]
	u.roundRobinQueue.index = (u.roundRobinQueue.index + 1) % length

	return element, nil
}

func (u Usecase) Remove(ctx context.Context, element string) {
	lock.Lock()
	defer lock.Unlock()

	for i, value := range u.roundRobinQueue.queue { //TODO: can optimize? but you dont usually have a large set of host so...
		if value == element {
			// Remove the element from the roundRobinQueue
			u.roundRobinQueue.queue = append(u.roundRobinQueue.queue[:i], u.roundRobinQueue.queue[i+1:]...)
			break
		}
	}
}
