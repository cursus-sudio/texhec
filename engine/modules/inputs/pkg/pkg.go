package inputspkg

import (
	"engine/modules/camera"
	"engine/modules/collider"
	engineinputs "engine/modules/inputs"
	"engine/modules/inputs/internal/mouse"
	"engine/modules/inputs/internal/systems"
	inputstool "engine/modules/inputs/internal/tool"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	"engine/services/media/inputs"
	"engine/services/media/window"
	"engine/services/runtime"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[engineinputs.Tool] { return inputstool.NewTool() })
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
					ioc.Get[logger.Logger](c),
					ioc.Get[ecs.ToolFactory[collider.CollisionTool]](c),
					ioc.Get[window.Api](c),
					ioc.Get[ecs.ToolFactory[camera.Tool]](c),
				),
				mouse.NewHoverSystem(),
				mouse.NewHoverEventsSystem(),
				mouse.NewClickSystem(ioc.Get[logger.Logger](c), ioc.Get[window.Api](c)),
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
