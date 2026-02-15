package constructpkg

import (
	"core/modules/construct"
	"core/modules/construct/internal"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) construct.Service {
		return internal.NewService(c)
	})
}
