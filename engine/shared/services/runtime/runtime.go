package runtime

import (
	"sync"
)

type Runtime interface {
	Run()
	Stop()
}

type runtime struct {
	onStart func()
	onStop  func()
	stopped bool
	mutex   sync.RWMutex
}

func NewRuntime(
	onStart func(),
	onStop func(),
) Runtime {
	r := &runtime{
		stopped: false,
		onStart: onStart,
		onStop:  onStop,
		mutex:   sync.RWMutex{},
	}
	r.mutex.Lock()
	return r
}

func (r *runtime) Run() {
	r.onStart()
	r.mutex.RLock()
	r.mutex.RUnlock()
}

func (r *runtime) Stop() {
	if r.stopped {
		return
	}
	r.onStop()
	r.mutex.Unlock()
	r.stopped = true
}
