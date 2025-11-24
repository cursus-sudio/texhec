package renderpkg

import (
	"frontend/modules/animation"
	"frontend/modules/render"
	"frontend/modules/render/internal"
	"frontend/services/media/window"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
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
				internal.NewErrorLogger(
					ioc.Get[logger.Logger](c),
					ioc.Get[ecs.ToolFactory[render.RenderTool]](c).Build(w),
				),
				internal.NewRenderSystem(ioc.Get[window.Api](c)),
			)
			return nil
		})
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b animation.AnimationSystemBuilder) animation.AnimationSystemBuilder {
		animation.AddTransitionFunction(b, func(w ecs.World) animation.TransitionFunction[render.ColorComponent] {
			componentArray := ecs.GetComponentsArray[render.ColorComponent](w)
			return func(arg animation.TransitionFunctionArgument[render.ColorComponent]) error {
				comp := arg.From.Blend(arg.To, float32(arg.State))
				return componentArray.SaveComponent(arg.Entity, comp)
			}
		})
		animation.AddTransitionFunction(b, func(w ecs.World) animation.TransitionFunction[render.TextureFrameComponent] {
			componentArray := ecs.GetComponentsArray[render.TextureFrameComponent](w)
			return func(arg animation.TransitionFunctionArgument[render.TextureFrameComponent]) error {
				comp := arg.From.Blend(arg.To, float64(arg.State))
				return componentArray.SaveComponent(arg.Entity, comp)
			}
		})
		return b
	})
}
