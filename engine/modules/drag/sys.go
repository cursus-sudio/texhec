package drag

import (
	"engine/modules/camera"
	"engine/modules/inputs"
	"engine/modules/transform"
	"engine/services/ecs"
)

type World interface {
	ecs.World
	transform.TransformTool
	camera.CameraTool
}

type System ecs.SystemRegister[World]

//

type DraggableEvent struct {
	Entity ecs.EntityID
	Drag   inputs.DragEvent
}

func NewDraggable(
	entity ecs.EntityID,
) DraggableEvent {
	return DraggableEvent{
		Entity: entity,
	}
}

func (e DraggableEvent) Apply(dragEvent inputs.DragEvent) any {
	e.Drag = dragEvent
	return e
}
