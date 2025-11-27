package scenes

import (
	"engine/services/ecs"
	"engine/services/scenes"
)

// change scene should happen after rendering
// because on scene change everything is cleaned up
type System ecs.SystemRegister

type ChangeSceneEvent struct {
	ID scenes.SceneId
}

func NewChangeSceneEvent(ID scenes.SceneId) ChangeSceneEvent {
	return ChangeSceneEvent{ID: ID}
}
