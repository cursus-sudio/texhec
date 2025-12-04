package dragpkg

import (
	"engine/modules/camera"
	"engine/modules/drag"
	"engine/modules/drag/internal"
	"engine/modules/transform"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			// events
			Register(drag.DraggableEvent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) drag.System {
		return internal.NewSystem(
			ioc.Get[logger.Logger](c),
			ioc.Get[ecs.ToolFactory[camera.Tool]](c),
			ioc.Get[ecs.ToolFactory[transform.Tool]](c),
		)
	})
}
