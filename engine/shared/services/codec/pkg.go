package codec

import (
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Builder {
		return NewBuilder()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) Codec {
		return ioc.Get[Builder](c).Build()
	})
	ioc.RegisterDependency[Codec, Builder](b)
}
