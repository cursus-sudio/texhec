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

func NewRenderSystem(
	world ecs.World,
	window window.Api,
	eventsBuilder events.Builder,
) render.System {
	return ecs.NewSystemRegister(func() error {
		s := &renderSystem{
			world:  world,
			events: eventsBuilder.Events(),
			window: window,
		}
		events.ListenE(eventsBuilder, s.Listen)
		return nil
	})
}

func (s *renderSystem) Listen(args frames.FrameEvent) error {
	events.Emit(s.events, render.RenderEvent{})

	s.window.Window().GLSwap()

	return nil
}
