package fpsloggerpkg

import (
	"core/modules/fpslogger"
	"core/modules/fpslogger/internal"
	"frontend/services/console"
	"frontend/services/scenes"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) fpslogger.System {
		return internal.NewFpsLoggerSystem(
			ioc.Get[scenes.SceneManager](c),
			ioc.Get[console.Console](c),
		)
	})
}
