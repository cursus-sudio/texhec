package runtime

import "github.com/ogiusek/ioc/v2"

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Builder {
		return newBuilder()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) Runtime {
		return ioc.Get[Builder](c).Build()
	})
}
