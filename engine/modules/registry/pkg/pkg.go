package registrypkg

import (
	"engine/modules/registry"
	"engine/modules/registry/internal"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) registry.Service {
		return internal.NewService(c)
	})
}
