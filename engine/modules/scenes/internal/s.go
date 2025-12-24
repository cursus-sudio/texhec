package internal

import (
	scenessys "engine/modules/scenes"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/scenes"

	"github.com/ogiusek/events"
)

type system struct {
	Manager scenes.SceneManager
	Event   *scenessys.ChangeSceneEvent
}

func NewChangeSceneSystem(m scenes.SceneManager) ecs.SystemRegister[ecs.World] {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &system{
			Manager: m,
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		events.ListenE(w.EventsBuilder(), s.ListenFrame)
		return nil
	})
}

func (s *system) Listen(event scenessys.ChangeSceneEvent) {
	s.Event = &event
}

func (s *system) ListenFrame(event frames.FrameEvent) error {
	if s.Event != nil {
		return s.Manager.LoadScene(s.Event.ID)
	}
	return nil
}
