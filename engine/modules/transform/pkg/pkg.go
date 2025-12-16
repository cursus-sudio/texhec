package transformpkg

import (
	"engine/modules/animation"
	"engine/modules/hierarchy"
	"engine/modules/transform"
	"engine/modules/transform/internal/transformtool"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	defaultPos         transform.PosComponent
	defaultRot         transform.RotationComponent
	defaultSize        transform.SizeComponent
	defaultPivot       transform.PivotPointComponent
	defaultParentPivot transform.ParentPivotPointComponent
}

func Package() ioc.Pkg {
	return pkg{
		transform.NewPos(0, 0, 0),
		transform.NewRotation(mgl32.QuatIdent()),
		transform.NewSize(1, 1, 1),
		transform.NewPivotPoint(.5, .5, .5),
		transform.NewParentPivotPoint(.5, .5, .5),
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			// components
			Register(transform.PosComponent{}).
			Register(transform.RotationComponent{}).
			Register(transform.SizeComponent{}).
			Register(transform.PivotPointComponent{}).
			Register(transform.ParentPivotPointComponent{}).
			Register(transform.ParentComponent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) transform.System {
		toolFactory := ioc.Get[ecs.ToolFactory[transform.Transform]](c)
		return ecs.NewSystemRegister(func(w ecs.World) error {
			toolFactory.Build(w)
			return nil
		})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[transform.Transform] {
		return transformtool.NewTransformTool(
			ioc.Get[logger.Logger](c),
			ioc.Get[ecs.ToolFactory[hierarchy.Hierarchy]](c),
			pkg.defaultRot,
			pkg.defaultSize,
			pkg.defaultPivot,
			pkg.defaultParentPivot,
		)
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b animation.AnimationSystemBuilder) animation.AnimationSystemBuilder {
		animation.AddTransitionFunction(b, func(w ecs.World) animation.TransitionFunction[transform.PosComponent] {
			componentArray := ecs.GetComponentsArray[transform.PosComponent](w)
			return func(arg animation.TransitionFunctionArgument[transform.PosComponent]) {
				comp := arg.From.Blend(arg.To, float32(arg.State))
				componentArray.SaveComponent(arg.Entity, comp)
			}
		})
		animation.AddTransitionFunction(b, func(w ecs.World) animation.TransitionFunction[transform.RotationComponent] {
			componentArray := ecs.GetComponentsArray[transform.RotationComponent](w)
			return func(arg animation.TransitionFunctionArgument[transform.RotationComponent]) {
				comp := arg.From.Blend(arg.To, float32(arg.State))
				componentArray.SaveComponent(arg.Entity, comp)
			}
		})
		animation.AddTransitionFunction(b, func(w ecs.World) animation.TransitionFunction[transform.SizeComponent] {
			componentArray := ecs.GetComponentsArray[transform.SizeComponent](w)
			return func(arg animation.TransitionFunctionArgument[transform.SizeComponent]) {
				comp := arg.From.Blend(arg.To, float32(arg.State))
				componentArray.SaveComponent(arg.Entity, comp)
			}
		})
		animation.AddTransitionFunction(b, func(w ecs.World) animation.TransitionFunction[transform.PivotPointComponent] {
			componentArray := ecs.GetComponentsArray[transform.PivotPointComponent](w)
			return func(arg animation.TransitionFunctionArgument[transform.PivotPointComponent]) {
				comp := arg.From.Blend(arg.To, float32(arg.State))
				componentArray.SaveComponent(arg.Entity, comp)
			}
		})
		animation.AddTransitionFunction(b, func(w ecs.World) animation.TransitionFunction[transform.ParentPivotPointComponent] {
			componentArray := ecs.GetComponentsArray[transform.ParentPivotPointComponent](w)
			return func(arg animation.TransitionFunctionArgument[transform.ParentPivotPointComponent]) {
				comp := arg.From.Blend(arg.To, float32(arg.State))
				componentArray.SaveComponent(arg.Entity, comp)
			}
		})
		return b
	})

}
