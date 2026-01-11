package drag

import (
	"engine/modules/inputs"
	"engine/services/ecs"
)

type System ecs.SystemRegister

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
