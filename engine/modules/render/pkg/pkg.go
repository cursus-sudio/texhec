package renderpkg

import (
	"engine/modules/animation"
	"engine/modules/render"
	"engine/modules/render/internal"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[render.Render] {
		return internal.NewTool()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.System {
		return ecs.NewSystemRegister(func(w ecs.World) error {
			ecs.RegisterSystems(w,
				internal.NewClearSystem(),
				internal.NewErrorLogger(
					ioc.Get[logger.Logger](c),
					ioc.Get[ecs.ToolFactory[render.Render]](c).Build(w),
				),
				internal.NewRenderSystem(ioc.Get[window.Api](c)),
			)
			return nil
		})
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b animation.AnimationSystemBuilder) animation.AnimationSystemBuilder {
		animation.AddTransitionFunction(b, func(w ecs.World) animation.TransitionFunction[render.ColorComponent] {
			componentArray := ecs.GetComponentsArray[render.ColorComponent](w)
			return func(arg animation.TransitionFunctionArgument[render.ColorComponent]) {
				comp := arg.From.Blend(arg.To, float32(arg.State))
				componentArray.SaveComponent(arg.Entity, comp)
			}
		})
		animation.AddTransitionFunction(b, func(w ecs.World) animation.TransitionFunction[render.TextureFrameComponent] {
			componentArray := ecs.GetComponentsArray[render.TextureFrameComponent](w)
			return func(arg animation.TransitionFunctionArgument[render.TextureFrameComponent]) {
				comp := arg.From.Blend(arg.To, float64(arg.State))
				componentArray.SaveComponent(arg.Entity, comp)
			}
		})
		return b
	})
}
