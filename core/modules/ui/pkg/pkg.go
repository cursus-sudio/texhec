package uipkg

import (
	"core/modules/tile"
	"core/modules/ui"
	"core/modules/ui/internal/uitool"
	"engine/modules/animation"
	"engine/modules/camera"
	"engine/modules/hierarchy"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"time"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	maxLayer          tile.Layer
	animationDuration time.Duration
	showAnimation     animation.AnimationID
	hideAnimation     animation.AnimationID
}

func Package(
	maxLayer tile.Layer,
	animationDuration time.Duration,
	showAnimation animation.AnimationID,
	hideAnimation animation.AnimationID,
) ioc.Pkg {
	return pkg{
		maxLayer,
		animationDuration,
		showAnimation,
		hideAnimation,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[ui.Tool] {
		return uitool.NewToolFactory(
			pkg.animationDuration,
			pkg.showAnimation,
			pkg.hideAnimation,
			ioc.Get[logger.Logger](c),
			ioc.Get[ecs.ToolFactory[camera.Tool]](c),
			ioc.Get[ecs.ToolFactory[transform.Tool]](c),
			ioc.Get[ecs.ToolFactory[tile.Tool]](c),
			ioc.Get[ecs.ToolFactory[text.Tool]](c),
			ioc.Get[ecs.ToolFactory[render.Tool]](c),
			ioc.Get[ecs.ToolFactory[hierarchy.Tool]](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) ui.System {
		factory := ioc.Get[ecs.ToolFactory[ui.Tool]](c)
		return ecs.NewSystemRegister(func(w ecs.World) error {
			factory.Build(w)
			return nil

		})
	})
	// ioc.RegisterSingleton(b, func(c ioc.Dic) ui.System {
	// 	return ecs.NewSystemRegister(func(w ecs.World) error {
	// 		errs := ecs.RegisterSystems(w,
	// 			uimodule.NewSystem(
	// 				ioc.Get[logger.Logger](c),
	// 				ioc.Get[ecs.ToolFactory[camera.Tool]](c),
	// 				ioc.Get[ecs.ToolFactory[transform.Tool]](c),
	// 				ioc.Get[ecs.ToolFactory[tile.Tool]](c),
	// 				ioc.Get[ecs.ToolFactory[text.Tool]](c),
	// 				ioc.Get[ecs.ToolFactory[render.Tool]](c),
	// 				ioc.Get[ecs.ToolFactory[hierarchy.Tool]](c),
	// 				pkg.maxLayer,
	// 			),
	// 			ecs.NewSystemRegister(func(w ecs.World) error {
	// 				events.Listen(w.EventsBuilder(), func(e sdl.MouseButtonEvent) {
	// 					if e.Button != sdl.BUTTON_RIGHT {
	// 						return
	// 					}
	// 					events.Emit(w.Events(), ui.HideUiEvent{})
	// 				})
	// 				return nil
	// 			}),
	// 		)
	// 		if len(errs) == 0 {
	// 			return nil
	// 		}
	// 		return errs[0]
	// 	})
	// })
}
