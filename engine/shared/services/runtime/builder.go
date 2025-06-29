package runtime

type Builder interface {
	OnStart(func())
	OnStop(func())
	Build() Runtime
}

type builder struct {
	onStart []func()
	onStop  []func()
}

func newBuilder() Builder {
	return &builder{}
}

func (b *builder) OnStart(listener func()) {
	b.onStart = append(b.onStart, listener)
}

func (b *builder) OnStop(listener func()) {
	b.onStop = append(b.onStop, listener)
}

func (b *builder) Build() Runtime {
	return NewRuntime(
		func() {
			for _, listner := range b.onStart {
				listner()
			}
		},
		func() {
			for _, listner := range b.onStop {
				listner()
			}
		},
	)
}
