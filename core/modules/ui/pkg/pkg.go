package uipkg

import (
	"core/modules/ui"
	"core/modules/ui/internal/uiservice"
	"engine/services/codec"
	"engine/services/ecs"
	"time"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type pkg struct {
	animationDuration time.Duration
}

func Package(
	animationDuration time.Duration,
) ioc.Pkg {
	return pkg{
		animationDuration,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// components
			Register(ui.UiCameraComponent{}).
			// events
			Register(ui.HideUiEvent{})
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ui.Service {
		return uiservice.NewService(c, pkg.animationDuration)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) ui.System {
		eventsBuilder := ioc.Get[events.Builder](c)
		return ecs.NewSystemRegister(func() error {
			events.Listen(eventsBuilder, func(e sdl.MouseButtonEvent) {
				if e.Button != sdl.BUTTON_RIGHT || e.State != sdl.RELEASED {
					return
				}
				events.Emit(eventsBuilder.Events(), ui.HideUiEvent{})
			})
			ioc.Get[ui.Service](c)
			return nil

		})
	})
}
