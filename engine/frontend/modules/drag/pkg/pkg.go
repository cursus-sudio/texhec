package dragpkg

import (
	"frontend/modules/camera"
	"frontend/modules/drag"
	"frontend/modules/drag/internal"
	"frontend/modules/transform"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) drag.System {
		return internal.NewSystem(
			ioc.Get[logger.Logger](c),
			ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
			ioc.Get[ecs.ToolFactory[transform.TransformTool]](c),
		)
	})
}
