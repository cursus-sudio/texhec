package inputspkg

import (
	"engine/modules/inputs"
	"engine/modules/inputs/internal/mouse"
	"engine/modules/inputs/internal/service"
	"engine/modules/inputs/internal/systems"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// components
			Register(inputs.HoveredComponent{}).
			Register(inputs.DraggedComponent{}).
			Register(inputs.KeepSelectedComponent{}).
			Register(inputs.LeftClickComponent{}).
			Register(inputs.DoubleLeftClickComponent{}).
			Register(inputs.RightClickComponent{}).
			Register(inputs.DoubleRightClickComponent{}).
			Register(inputs.MouseEnterComponent{}).
			Register(inputs.MouseLeaveComponent{}).
			Register(inputs.HoverComponent{}).
			Register(inputs.DragComponent{}).
			// events
			Register(inputs.QuitEvent{}).
			Register(inputs.DragEvent{}).
			Register(inputs.SynchronizePositionEvent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) inputs.Service {
		return service.NewService(c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) inputs.System {
		return ecs.NewSystemRegister(func() error {
			ecs.RegisterSystems(
				systems.NewInputsSystem(c),

				ecs.NewSystemRegister(func() error {
					eventsBuilder := ioc.Get[events.Builder](c)
					events.Listen(eventsBuilder, func(sdl.QuitEvent) {
						events.Emit(eventsBuilder.Events(), inputs.NewQuitEvent())
					})
					return nil
				}),

				mouse.NewCameraRaySystem(c),
				mouse.NewHoverSystem(c),
				mouse.NewHoverEventsSystem(c),
				mouse.NewClickSystem(c),
				ecs.NewSystemRegister(func() error {
					eventsBuilder := ioc.Get[events.Builder](c)
					events.Listen(eventsBuilder, func(frames.FrameEvent) {
						events.Emit(eventsBuilder.Events(), mouse.NewShootRayEvent())
					})
					return nil
				}),
			)
			return nil
		})
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) inputs.ShutdownSystem {
		return systems.NewQuitSystem(c)
	})
}
