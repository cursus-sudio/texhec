package inputs

import (
	"frontend/services/frames"
	"shared/services/clock"
	"shared/services/logger"
	"shared/services/runtime"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) *api {
		return newInputsApi(
			ioc.Get[logger.Logger](c),
			ioc.Get[clock.Clock](c),
			ioc.Get[events.Events](c),
		)
	})
	ioc.RegisterDependency[*api, logger.Logger](b)
	ioc.RegisterDependency[*api, clock.Clock](b)
	ioc.RegisterDependency[*api, events.Events](b)

	ioc.RegisterTransient(b, func(c ioc.Dic) Api {
		return ioc.Get[*api](c)
	})
	ioc.RegisterDependency[Api, *api](b)

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b events.Builder) events.Builder {
		events.Listen(b, func(qe sdl.QuitEvent) {
			ioc.Get[runtime.Runtime](c).Stop()
		})
		return b
	})

	ioc.WrapService(b, frames.HandleInputs, func(c ioc.Dic, b events.Builder) events.Builder {
		var i *api
		events.Listen(b, func(e frames.FrameEvent) {
			if i == nil {
				i = ioc.Get[*api](c)
			}
			i.Poll()
		})
		return b
	})
}
