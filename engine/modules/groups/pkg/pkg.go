package groupspkg

import (
	"engine/modules/groups"
	"engine/modules/groups/internal"
	"engine/modules/hierarchy"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) groups.System {
		return internal.NewSystem(
			ioc.Get[logger.Logger](c),
			ioc.Get[ecs.ToolFactory[hierarchy.Tool]](c),
		)
	})
}
