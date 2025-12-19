package frames

import (
	"engine/services/clock"
	runtimeservice "engine/services/runtime"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	tps,
	fps int
}

func Package(tps, fps int) ioc.Pkg {
	return pkg{
		tps: tps,
		fps: fps,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Builder {
		return NewBuilder(pkg.tps, pkg.fps)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) Frames {
		return ioc.Get[Builder](c).Build(ioc.Get[events.Events](c), ioc.Get[clock.Clock](c))
	})
	ioc.RegisterDependency[Frames, events.Events](b)
	ioc.RegisterDependency[Frames, clock.Clock](b)
	ioc.RegisterDependency[Frames, Builder](b)

	ioc.WrapService(b, runtimeservice.OrderStop, func(c ioc.Dic, r runtimeservice.Builder) runtimeservice.Builder {
		r.OnStartOnMainThread(func(r runtimeservice.Runtime) {
			ioc.Get[Frames](c).Run()
			go r.Stop()
		})
		r.OnStop(func(r runtimeservice.Runtime) {
			ioc.Get[Frames](c).Stop()
		})
		return r
	})
	ioc.RegisterDependency[runtimeservice.Builder, Frames](b)
}
