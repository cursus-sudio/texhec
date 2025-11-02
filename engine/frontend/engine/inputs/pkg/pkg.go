package inputspkg

import (
	"frontend/engine/camera"
	"frontend/engine/collider"
	engineinputs "frontend/engine/inputs"
	"frontend/engine/inputs/internal/mouse"
	"frontend/engine/inputs/internal/systems"
	"frontend/services/frames"
	"frontend/services/media/inputs"
	"frontend/services/media/window"
	"shared/services/ecs"
	"shared/services/logger"
	"shared/services/runtime"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) engineinputs.System {
		return ecs.NewSystemRegister(func(w ecs.World) error {
			ecs.RegisterSystems(w,
				systems.NewInputsSystem(ioc.Get[inputs.Api](c)),
				systems.NewResizeSystem(),
				systems.NewQuitSystem(ioc.Get[runtime.Runtime](c)),

				ecs.NewSystemRegister(func(w ecs.World) error {
					events.Listen(w.EventsBuilder(), func(sdl.QuitEvent) {
						events.Emit(w.Events(), engineinputs.NewQuitEvent())
					})
					return nil
				}),

				mouse.NewCameraRaySystem(
					ioc.Get[ecs.ToolFactory[collider.CollisionTool]](c),
					ioc.Get[window.Api](c),
					ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
				),
				mouse.NewHoverSystem(),
				mouse.NewHoverEventsSystem(),
				mouse.NewClickSystem(ioc.Get[logger.Logger](c)),
				ecs.NewSystemRegister(func(w ecs.World) error {
					events.Listen(w.EventsBuilder(), func(frames.FrameEvent) {
						events.Emit(w.Events(), mouse.NewShootRayEvent())
					})
					return nil
				}),
			)
			return nil
		})
	})
}
