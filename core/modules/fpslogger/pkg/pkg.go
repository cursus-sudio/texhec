package fpsloggerpkg

import (
	"core/modules/fpslogger"
	"core/modules/fpslogger/internal"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) fpslogger.System {
		return internal.NewFpsLoggerSystem(c)
	})
}
