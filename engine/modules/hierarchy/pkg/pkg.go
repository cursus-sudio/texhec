package hierarchypkg

import (
	"engine/modules/hierarchy"
	"engine/modules/hierarchy/internal/hierarchytool"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[hierarchy.Tool] {
		return hierarchytool.NewTool(
			ioc.Get[logger.Logger](c),
		)
	})
}
