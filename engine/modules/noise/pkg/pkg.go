package noisepkg

import (
	"engine/modules/noise"
	"engine/modules/noise/internal"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) noise.Service {
		return internal.NewService(c)
	})
}
