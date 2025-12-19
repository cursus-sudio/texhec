package slerppkg

import (
	"engine/modules/slerp"
	"engine/modules/slerp/internal/sys"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) sys.Builder {
		return sys.NewBuilder()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) slerp.System {
		return ioc.Get[sys.Builder](c).Build()
	})
}
