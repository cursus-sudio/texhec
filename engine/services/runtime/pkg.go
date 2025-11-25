package runtime

import (
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Builder {
		return newBuilder()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) Runtime {
		return ioc.Get[Builder](c).Build()
	})
	ioc.RegisterDependency[Runtime, Builder](b)
}
