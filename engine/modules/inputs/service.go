package inputs

import (
	"engine/modules/collider"
	"engine/services/ecs"
)

type Service interface {
	Hovered() ecs.ComponentsArray[HoveredComponent]
	Dragged() ecs.ComponentsArray[DraggedComponent]
	Stacked() ecs.ComponentsArray[StackedComponent]

	KeepSelected() ecs.ComponentsArray[KeepSelectedComponent]

	LeftClick() ecs.ComponentsArray[LeftClickComponent]
	DoubleLeftClick() ecs.ComponentsArray[DoubleLeftClickComponent]

	RightClick() ecs.ComponentsArray[RightClickComponent]
	DoubleRightClick() ecs.ComponentsArray[DoubleRightClickComponent]

	MouseEnter() ecs.ComponentsArray[MouseEnterComponent]
	MouseLeave() ecs.ComponentsArray[MouseLeaveComponent]

	Hover() ecs.ComponentsArray[HoverComponent]
	Drag() ecs.ComponentsArray[DragComponent]

	Stack() ecs.ComponentsArray[StackComponent]

	// returns ordered targets with additional data
	StackedData() []Target
}

type EventTargetSetter interface {
	SetTarget(Target) EventTargetSetter
}

type Target struct {
	collider.ObjectRayCollision
	Camera ecs.EntityID
}
