package codec

import "github.com/ogiusek/ioc/v2"

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Builder {
		return NewBuilder()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) Codec {
		return ioc.Get[Builder](c).Build()
	})
}
