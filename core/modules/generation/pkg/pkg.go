package generationpkg

import (
	"core/modules/generation"
	"core/modules/generation/internal"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) generation.Service {
		return internal.NewService(c)
	})
}
