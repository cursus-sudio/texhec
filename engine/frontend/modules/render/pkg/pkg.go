package renderpkg

import (
	"frontend/modules/render"
	"frontend/modules/render/internal"
	"frontend/services/media/window"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[render.RenderTool] {
		return ecs.NewToolFactory(func(w ecs.World) render.RenderTool {
			return internal.NewTool()
		})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.RenderTool {
		return internal.NewTool()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.System {
		return ecs.NewSystemRegister(func(w ecs.World) error {
			ecs.RegisterSystems(w,
				internal.NewClearSystem(),
				// TODO expose error cause
				internal.NewErrorLogger(
					ioc.Get[logger.Logger](c),
					ioc.Get[ecs.ToolFactory[render.RenderTool]](c).Build(w),
				),
				internal.NewRenderSystem(ioc.Get[window.Api](c)),
			)
			return nil
		})
	})
}
