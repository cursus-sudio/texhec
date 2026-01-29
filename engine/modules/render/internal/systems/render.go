package systems

import (
	"engine/modules/render"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/media/window"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type renderSystem struct {
	World         ecs.World      `inject:"1"`
	Events        events.Events  `inject:"1"`
	Window        window.Api     `inject:"1"`
	EventsBuilder events.Builder `inject:"1"`
}

func NewRenderSystem(c ioc.Dic) render.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*renderSystem](c)
		events.ListenE(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *renderSystem) Listen(args frames.FrameEvent) error {
	events.Emit(s.Events, render.RenderEvent{})

	s.Window.Window().GLSwap()

	return nil
}
