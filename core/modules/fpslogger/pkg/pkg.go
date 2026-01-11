package fpsloggerpkg

import (
	"core/modules/fpslogger"
	"core/modules/fpslogger/internal"
	"engine/services/console"
	"engine/services/ecs"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) fpslogger.System {
		return internal.NewFpsLoggerSystem(
			ioc.Get[events.Builder](c),
			ioc.Get[ecs.World](c),
			ioc.Get[console.Console](c),
		)
	})
}
