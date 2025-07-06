package frames

import (
	"shared/services/clock"
	runtimeservice "shared/services/runtime"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Builder {
		return NewBuilder(ioc.Get[clock.Clock](c))
	})
	ioc.RegisterDependency[Builder, clock.Clock](b)

	ioc.RegisterSingleton(b, func(c ioc.Dic) Frames {
		return ioc.Get[Builder](c).Build()
	})
	ioc.RegisterDependency[Frames, Builder](b)

	ioc.WrapService(b, runtimeservice.OrderStop, func(c ioc.Dic, r runtimeservice.Builder) runtimeservice.Builder {
		r.OnStartOnMainThread(func(r runtimeservice.Runtime) {
			ioc.Get[Frames](c).Run()
			r.Stop()
		})
		r.OnStop(func(r runtimeservice.Runtime) {
			ioc.Get[Frames](c).Stop()
		})
		return r
	})
	ioc.RegisterDependency[runtimeservice.Builder, Frames](b)
}
