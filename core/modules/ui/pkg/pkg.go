package uipkg

import (
	gameassets "core/assets"
	"core/modules/tile"
	"core/modules/ui"
	"core/modules/ui/internal/uitool"
	"engine/modules/animation"
	"engine/modules/camera"
	"engine/modules/hierarchy"
	"engine/modules/render"
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
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
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			// components
			Register(ui.UiCameraComponent{}).
			// events
			Register(ui.HideUiEvent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[ui.UiTool] {
		return uitool.NewToolFactory(
			pkg.animationDuration,
			pkg.showAnimation,
			pkg.hideAnimation,
			ioc.Get[gameassets.GameAssets](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
			ioc.Get[ecs.ToolFactory[transform.TransformTool]](c),
			ioc.Get[ecs.ToolFactory[tile.Tile]](c),
			ioc.Get[ecs.ToolFactory[text.TextTool]](c),
			ioc.Get[ecs.ToolFactory[render.RenderTool]](c),
			ioc.Get[ecs.ToolFactory[hierarchy.HierarchyTool]](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) ui.System {
		factory := ioc.Get[ecs.ToolFactory[ui.UiTool]](c)
		return ecs.NewSystemRegister(func(w ecs.World) error {
			events.Listen(w.EventsBuilder(), func(e sdl.MouseButtonEvent) {
				if e.Button != sdl.BUTTON_RIGHT || e.State != sdl.RELEASED {
					return
				}
				events.Emit(w.Events(), ui.HideUiEvent{})
			})
			factory.Build(w)
			return nil

		})
	})
}
