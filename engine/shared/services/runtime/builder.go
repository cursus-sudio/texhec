package runtime

import (
	"shared/utils/httperrors"

	"github.com/ogiusek/ioc/v2"
)

const (
	OrderStop ioc.Order = iota
	OrderCleanUp
)

type Builder interface {
	BeforeStart(func(Runtime))
	OnStart(func(Runtime))
	OnStartOnMainThread(func(Runtime)) error
	OnStop(func(Runtime))
	Build() Runtime
}

type builder struct {
	beforeStart         []func(Runtime)
	onStart             []func(Runtime)
	onStartOnMainThread func(Runtime)
	onStop              []func(Runtime)
}

func newBuilder() Builder {
	return &builder{}
}

func (b *builder) BeforeStart(listener func(Runtime)) {
	b.beforeStart = append(b.beforeStart, listener)
}

func (b *builder) OnStart(listener func(Runtime)) {
	b.onStart = append(b.onStart, listener)
}

func (b *builder) OnStartOnMainThread(listener func(Runtime)) error {
	if b.onStartOnMainThread != nil {
		return httperrors.Err409
	}
	b.onStartOnMainThread = listener
	return nil
}

func (b *builder) OnStop(listener func(Runtime)) {
	b.onStop = append(b.onStop, listener)
}

func (b *builder) Build() Runtime {
	return NewRuntime(
		func(r Runtime) {
			for _, listener := range b.beforeStart {
				listener(r)
			}
			for _, listener := range b.onStart {
				go listener(r)
			}
			if b.onStartOnMainThread != nil {
				b.onStartOnMainThread(r)
			}
		},
		func(r Runtime) {
			for _, listner := range b.onStop {
				listner(r)
			}
		},
	)
}
