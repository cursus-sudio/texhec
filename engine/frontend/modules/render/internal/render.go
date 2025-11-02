package internal

import (
	"frontend/modules/render"
	"frontend/services/frames"
	"frontend/services/media/window"
	"shared/services/ecs"
	"sync"

	"github.com/ogiusek/events"
)

type renderSystem struct {
	world  ecs.World
	events events.Events
	window window.Api

	mutex sync.Locker
}

func NewRenderSystem(window window.Api) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &renderSystem{
			world:  w,
			events: w.Events(),
			window: window,

			mutex: &sync.Mutex{},
		}
		events.ListenE(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *renderSystem) Listen(args frames.FrameEvent) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	events.Emit(s.events, render.RenderEvent{})

	s.window.Window().GLSwap()

	return nil
}
