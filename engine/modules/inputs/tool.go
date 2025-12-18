package inputs

import (
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/services/ecs"
)

type InputsTool interface {
	Inputs() Interface
}

type World interface {
	ecs.World
	collider.ColliderTool
	camera.CameraTool
}

type Interface interface {
	Hovered() ecs.ComponentsArray[HoveredComponent]
	Dragged() ecs.ComponentsArray[DraggedComponent]

	KeepSelected() ecs.ComponentsArray[KeepSelectedComponent]

	MouseLeft() ecs.ComponentsArray[MouseLeftClickComponent]
	MouseDoubleLeft() ecs.ComponentsArray[MouseDoubleLeftClickComponent]

	MouseRight() ecs.ComponentsArray[MouseRightClickComponent]
	MouseDoubleRight() ecs.ComponentsArray[MouseDoubleRightClickComponent]

	MouseEnter() ecs.ComponentsArray[MouseEnterComponent]
	MouseLeave() ecs.ComponentsArray[MouseLeaveComponent]

	MouseHover() ecs.ComponentsArray[MouseHoverComponent]
	MouseDrag() ecs.ComponentsArray[MouseDragComponent]
}
