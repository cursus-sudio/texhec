package inputs

import (
	"engine/services/clock"
	"engine/services/logger"
	"engine/services/runtime"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Api {
		return newInputsApi(
			ioc.Get[logger.Logger](c),
			ioc.Get[clock.Clock](c),
			ioc.Get[events.Events](c),
		)
	})
	ioc.RegisterDependency[Api, logger.Logger](b)
	ioc.RegisterDependency[Api, clock.Clock](b)
	ioc.RegisterDependency[Api, events.Events](b)

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b events.Builder) events.Builder {
		events.Listen(b, func(qe sdl.QuitEvent) {
			ioc.Get[runtime.Runtime](c).Stop()
		})
		return b
	})
}
