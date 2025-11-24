package uipkg

import (
	"core/modules/tile"
	"core/modules/ui"
	"core/modules/ui/internal/uimodule"
	"frontend/modules/camera"
	"frontend/modules/text"
	"frontend/modules/transform"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type pkg struct {
	maxLayer tile.Layer
}

func Package(
	maxLayer tile.Layer,
) ioc.Pkg {
	return pkg{maxLayer}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ui.System {
		return ecs.NewSystemRegister(func(w ecs.World) error {
			errs := ecs.RegisterSystems(w,
				uimodule.NewSystem(
					ioc.Get[logger.Logger](c),
					ioc.Get[ecs.ToolFactory[camera.Tool]](c),
					ioc.Get[ecs.ToolFactory[transform.Tool]](c),
					ioc.Get[ecs.ToolFactory[tile.Tool]](c),
					ioc.Get[ecs.ToolFactory[text.Tool]](c),
					pkg.maxLayer,
				),
				ecs.NewSystemRegister(func(w ecs.World) error {
					events.Listen(w.EventsBuilder(), func(e sdl.MouseButtonEvent) {
						if e.Button != sdl.BUTTON_RIGHT {
							return
						}
						events.Emit(w.Events(), ui.UnselectEvent{})
					})
					return nil
				}),
			)
			if len(errs) == 0 {
				return nil
			}
			return errs[0]
		})
	})
}
