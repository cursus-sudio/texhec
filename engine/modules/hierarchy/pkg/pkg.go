package hierarchypkg

import (
	"engine/modules/hierarchy"
	"engine/modules/hierarchy/internal/hierarchytool"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			// components
			Register(hierarchy.ParentComponent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[hierarchy.Tool] {
		return hierarchytool.NewTool(
			ioc.Get[logger.Logger](c),
		)
	})
}
