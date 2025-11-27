package inputs

import "engine/services/ecs"

type Tool interface {
	Transaction() Transaction
}

type Transaction interface {
	GetObject(ecs.EntityID) Object
	Transactions() []ecs.AnyComponentsArrayTransaction
	Flush() error
}

type Object interface {
	Hovered() ecs.EntityComponent[HoveredComponent]
	Dragged() ecs.EntityComponent[DraggedComponent]

	KeepSelected() ecs.EntityComponent[KeepSelectedComponent]

	MouseLeftClick() ecs.EntityComponent[MouseLeftClickComponent]
	MouseDoubleLeftClick() ecs.EntityComponent[MouseDoubleLeftClickComponent]

	MouseRightClick() ecs.EntityComponent[MouseRightClickComponent]
	MouseDoubleRightClick() ecs.EntityComponent[MouseDoubleRightClickComponent]

	MouseEnter() ecs.EntityComponent[MouseEnterComponent]
	MouseLeave() ecs.EntityComponent[MouseLeaveComponent]

	MouseHover() ecs.EntityComponent[MouseHoverComponent]
	MouseDrag() ecs.EntityComponent[MouseDragComponent]
}
