package scenessys

import (
	"frontend/services/scenes"
	"shared/services/ecs"

	"github.com/ogiusek/events"
)

type ChangeSceneEvent struct {
	ID scenes.SceneId
}

func NewChangeSceneEvent(ID scenes.SceneId) ChangeSceneEvent {
	return ChangeSceneEvent{ID: ID}
}

type system struct {
	Manager scenes.SceneManager
}

func NewChangeSceneSystem(m scenes.SceneManager) ecs.SystemRegister {
	return &system{
		Manager: m,
	}
}

func (s *system) Register(b events.Builder) {
	events.ListenE(b, s.Listen)
}

func (s *system) Listen(event ChangeSceneEvent) error {
	return s.Manager.LoadScene(event.ID)
}
