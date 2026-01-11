package scene

import (
	"engine/services/ecs"
)

// change scene should happen after rendering
// because on scene change everything is cleaned up
type System ecs.SystemRegister[ecs.World]

type ChangeSceneEvent struct {
	ID ID
}

func NewChangeSceneEvent(ID ID) ChangeSceneEvent {
	return ChangeSceneEvent{ID: ID}
}

//

type ID struct {
	ID string
}

func NewSceneId(id string) ID {
	return ID{id}
}

//

type Scene func(sceneParent ecs.EntityID)

type Service interface {
	SetScene(ID, Scene)
}
