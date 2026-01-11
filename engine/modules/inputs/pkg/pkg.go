package inputspkg

import (
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/inputs"
	"engine/modules/inputs/internal/mouse"
	"engine/modules/inputs/internal/service"
	"engine/modules/inputs/internal/systems"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	inputsapi "engine/services/media/inputs"
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
		return service.NewService(
			ioc.Get[logger.Logger](c),
			ioc.Get[events.Builder](c),
			ioc.Get[ecs.World](c),
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) inputs.System {
		return ecs.NewSystemRegister(func() error {
			ecs.RegisterSystems(
				systems.NewInputsSystem(
					ioc.Get[events.Builder](c),
					ioc.Get[inputsapi.Api](c),
				),
				systems.NewResizeSystem(
					ioc.Get[events.Builder](c),
				),
				systems.NewQuitSystem(
					ioc.Get[runtime.Runtime](c),
					ioc.Get[events.Builder](c),
				),

				ecs.NewSystemRegister(func() error {
					eventsBuilder := ioc.Get[events.Builder](c)
					events.Listen(eventsBuilder, func(sdl.QuitEvent) {
						events.Emit(eventsBuilder.Events(), inputs.NewQuitEvent())
					})
					return nil
				}),

				mouse.NewCameraRaySystem(
					ioc.Get[events.Builder](c),
					ioc.Get[ecs.World](c),
					ioc.Get[camera.Service](c),
					ioc.Get[collider.Service](c),
					ioc.Get[logger.Logger](c),
					ioc.Get[window.Api](c),
				),
				mouse.NewHoverSystem(
					ioc.Get[events.Builder](c),
					ioc.Get[ecs.World](c),
					ioc.Get[inputs.Service](c),
					ioc.Get[logger.Logger](c),
				),
				mouse.NewHoverEventsSystem(
					ioc.Get[events.Builder](c),
					ioc.Get[ecs.World](c),
					ioc.Get[inputs.Service](c),
				),
				mouse.NewClickSystem(
					ioc.Get[logger.Logger](c),
					ioc.Get[window.Api](c),
					ioc.Get[events.Builder](c),
					ioc.Get[ecs.World](c),
					ioc.Get[inputs.Service](c),
				),
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
}
