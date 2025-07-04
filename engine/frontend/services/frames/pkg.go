package frames

import (
	"shared/services/clock"
	"shared/services/runtime"

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

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, r runtime.Builder) runtime.Builder {
		r.OnStart(func() {
			go ioc.Get[Frames](c).Run()
		})
		return r
	})
	ioc.RegisterDependency[runtime.Builder, Frames](b)
}
