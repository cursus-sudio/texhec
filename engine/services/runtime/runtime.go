package runtime

import (
	"sync"
)

type Runtime interface {
	Run()
	Stop()
}

type runtime struct {
	onStart   func(Runtime)
	onStop    func(Runtime)
	stopped   bool
	stopMutex sync.Mutex
	mutex     sync.RWMutex
}

func NewRuntime(
	onStart func(Runtime),
	onStop func(Runtime),
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
	r.onStart(r)
	r.mutex.RLock()
}

func (r *runtime) Stop() {
	// panic("went")
	r.stopMutex.Lock()
	if r.stopped {
		return
	}
	r.stopped = true
	r.stopMutex.Unlock()
	go func() {
		r.onStop(r)
		r.mutex.Unlock()
	}()
}
