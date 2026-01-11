package hierarchypkg

import (
	"engine/modules/hierarchy"
	"engine/modules/hierarchy/internal/hierarchyservice"
	"engine/services/codec"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// components
			Register(hierarchy.Component{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) hierarchy.Service {
		return hierarchyservice.NewService(c)
	})
}
