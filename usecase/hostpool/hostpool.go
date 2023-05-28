package hostpool

import (
	"context"
	"errors"
	"github.com/IshlahulHanif/poneglyph"
)

func (u *Usecase) GetHost(ctx context.Context) (res string, err error) {
	lock.Lock()
	defer lock.Unlock()

	// get current host in queue
	length := len(u.roundRobinQueue.queue)
	if length == 0 {
		// something wrong
		return "", poneglyph.Trace(errors.New("hostpool empty"))
	}

	if length-1 < u.roundRobinQueue.index {
		// index overload
		u.roundRobinQueue.index = 0
		// TODO: move host around to avoid overload on first host
	}

	res = u.roundRobinQueue.queue[u.roundRobinQueue.index]
	u.roundRobinQueue.index = (u.roundRobinQueue.index + 1) % length

	return res, nil
}

func (u *Usecase) AddHost(ctx context.Context, host string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	u.roundRobinQueue.queue = append(u.roundRobinQueue.queue, host)

	return nil
}

func (u *Usecase) RemoveHost(ctx context.Context, host string) (err error) {
	lock.Lock()
	defer lock.Unlock()

	for i, value := range u.roundRobinQueue.queue { //TODO: can optimize? but you dont usually have a large set of host so...
		if value == host {
			// Remove the element from the roundRobinQueue
			u.roundRobinQueue.queue = append(u.roundRobinQueue.queue[:i], u.roundRobinQueue.queue[i+1:]...)
			// no need to move index
			break
		}
	}

	return nil
}
