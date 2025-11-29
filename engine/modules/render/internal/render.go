package internal

import (
	"engine/modules/render"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/media/window"

	"github.com/ogiusek/events"
)

type renderSystem struct {
	world  ecs.World
	events events.Events
	window window.Api
}

func NewRenderSystem(window window.Api) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &renderSystem{
			world:  w,
			events: w.Events(),
			window: window,
		}
		events.ListenE(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *renderSystem) Listen(args frames.FrameEvent) error {
	events.Emit(s.events, render.RenderEvent{})

	s.window.Window().GLSwap()

	return nil
}
