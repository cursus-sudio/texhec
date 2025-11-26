package inputstool

import (
	"engine/modules/inputs"
	"engine/services/ecs"
)

type tool struct {
	world ecs.World

	hoveredArray               ecs.ComponentsArray[inputs.HoveredComponent]
	draggedArray               ecs.ComponentsArray[inputs.DraggedComponent]
	keepSelectedArray          ecs.ComponentsArray[inputs.KeepSelectedComponent]
	mouseLeftClickArray        ecs.ComponentsArray[inputs.MouseLeftClickComponent]
	mouseDoubleLeftClickArray  ecs.ComponentsArray[inputs.MouseDoubleLeftClickComponent]
	mouseRightClickArray       ecs.ComponentsArray[inputs.MouseRightClickComponent]
	mouseDoubleRightClickArray ecs.ComponentsArray[inputs.MouseDoubleRightClickComponent]
	mouseEnterArray            ecs.ComponentsArray[inputs.MouseEnterComponent]
	mouseLeaveArray            ecs.ComponentsArray[inputs.MouseLeaveComponent]
	mouseHoverArray            ecs.ComponentsArray[inputs.MouseHoverComponent]
	mouseDragArray             ecs.ComponentsArray[inputs.MouseDragComponent]
}

func NewTool() ecs.ToolFactory[inputs.Tool] {
	return ecs.NewToolFactory(func(w ecs.World) inputs.Tool {
		return &tool{
			w,
			ecs.GetComponentsArray[inputs.HoveredComponent](w),
			ecs.GetComponentsArray[inputs.DraggedComponent](w),
			ecs.GetComponentsArray[inputs.KeepSelectedComponent](w),
			ecs.GetComponentsArray[inputs.MouseLeftClickComponent](w),
			ecs.GetComponentsArray[inputs.MouseDoubleLeftClickComponent](w),
			ecs.GetComponentsArray[inputs.MouseRightClickComponent](w),
			ecs.GetComponentsArray[inputs.MouseDoubleRightClickComponent](w),
			ecs.GetComponentsArray[inputs.MouseEnterComponent](w),
			ecs.GetComponentsArray[inputs.MouseLeaveComponent](w),
			ecs.GetComponentsArray[inputs.MouseHoverComponent](w),
			ecs.GetComponentsArray[inputs.MouseDragComponent](w),
		}
	})
}

//

type transaction struct {
	*tool

	hoveredTransaction               ecs.ComponentsArrayTransaction[inputs.HoveredComponent]
	draggedTransaction               ecs.ComponentsArrayTransaction[inputs.DraggedComponent]
	keepSelectedTransaction          ecs.ComponentsArrayTransaction[inputs.KeepSelectedComponent]
	mouseLeftClickTransaction        ecs.ComponentsArrayTransaction[inputs.MouseLeftClickComponent]
	mouseDoubleLeftClickTransaction  ecs.ComponentsArrayTransaction[inputs.MouseDoubleLeftClickComponent]
	mouseRightClickTransaction       ecs.ComponentsArrayTransaction[inputs.MouseRightClickComponent]
	mouseDoubleRightClickTransaction ecs.ComponentsArrayTransaction[inputs.MouseDoubleRightClickComponent]
	mouseEnterTransaction            ecs.ComponentsArrayTransaction[inputs.MouseEnterComponent]
	mouseLeaveTransaction            ecs.ComponentsArrayTransaction[inputs.MouseLeaveComponent]
	mouseHoverTransaction            ecs.ComponentsArrayTransaction[inputs.MouseHoverComponent]
	mouseDragTransaction             ecs.ComponentsArrayTransaction[inputs.MouseDragComponent]
}

func (t *tool) Transaction() inputs.Transaction {
	return &transaction{
		t,
		t.hoveredArray.Transaction(),
		t.draggedArray.Transaction(),
		t.keepSelectedArray.Transaction(),
		t.mouseLeftClickArray.Transaction(),
		t.mouseDoubleLeftClickArray.Transaction(),
		t.mouseRightClickArray.Transaction(),
		t.mouseDoubleRightClickArray.Transaction(),
		t.mouseEnterArray.Transaction(),
		t.mouseLeaveArray.Transaction(),
		t.mouseHoverArray.Transaction(),
		t.mouseDragArray.Transaction(),
	}
}

type object struct {
	*transaction

	hovered               ecs.EntityComponent[inputs.HoveredComponent]
	dragged               ecs.EntityComponent[inputs.DraggedComponent]
	keepSelected          ecs.EntityComponent[inputs.KeepSelectedComponent]
	mouseLeftClick        ecs.EntityComponent[inputs.MouseLeftClickComponent]
	mouseDoubleLeftClick  ecs.EntityComponent[inputs.MouseDoubleLeftClickComponent]
	mouseRightClick       ecs.EntityComponent[inputs.MouseRightClickComponent]
	mouseDoubleRightClick ecs.EntityComponent[inputs.MouseDoubleRightClickComponent]
	mouseEnter            ecs.EntityComponent[inputs.MouseEnterComponent]
	mouseLeave            ecs.EntityComponent[inputs.MouseLeaveComponent]
	mouseHover            ecs.EntityComponent[inputs.MouseHoverComponent]
	mouseDrag             ecs.EntityComponent[inputs.MouseDragComponent]
}

func (t *transaction) GetObject(entity ecs.EntityID) inputs.Object {
	return &object{
		t,
		t.hoveredTransaction.GetEntityComponent(entity),
		t.draggedTransaction.GetEntityComponent(entity),
		t.keepSelectedTransaction.GetEntityComponent(entity),
		t.mouseLeftClickTransaction.GetEntityComponent(entity),
		t.mouseDoubleLeftClickTransaction.GetEntityComponent(entity),
		t.mouseRightClickTransaction.GetEntityComponent(entity),
		t.mouseDoubleRightClickTransaction.GetEntityComponent(entity),
		t.mouseEnterTransaction.GetEntityComponent(entity),
		t.mouseLeaveTransaction.GetEntityComponent(entity),
		t.mouseHoverTransaction.GetEntityComponent(entity),
		t.mouseDragTransaction.GetEntityComponent(entity),
	}
}
func (t *transaction) Transactions() []ecs.AnyComponentsArrayTransaction {
	return []ecs.AnyComponentsArrayTransaction{
		t.hoveredTransaction,
		t.draggedTransaction,
		t.keepSelectedTransaction,
		t.mouseLeftClickTransaction,
		t.mouseDoubleLeftClickTransaction,
		t.mouseRightClickTransaction,
		t.mouseDoubleRightClickTransaction,
		t.mouseEnterTransaction,
		t.mouseHoverTransaction,
		t.mouseDragTransaction,
	}
}
func (t *transaction) Flush() error { return ecs.FlushMany(t.Transactions()...) }

//

func (t *object) Hovered() ecs.EntityComponent[inputs.HoveredComponent] { return t.hovered }
func (t *object) Dragged() ecs.EntityComponent[inputs.DraggedComponent] { return t.dragged }
func (t *object) KeepSelected() ecs.EntityComponent[inputs.KeepSelectedComponent] {
	return t.keepSelected
}
func (t *object) MouseLeftClick() ecs.EntityComponent[inputs.MouseLeftClickComponent] {
	return t.mouseLeftClick
}
func (t *object) MouseDoubleLeftClick() ecs.EntityComponent[inputs.MouseDoubleLeftClickComponent] {
	return t.mouseDoubleLeftClick
}
func (t *object) MouseRightClick() ecs.EntityComponent[inputs.MouseRightClickComponent] {
	return t.mouseRightClick
}
func (t *object) MouseDoubleRightClick() ecs.EntityComponent[inputs.MouseDoubleRightClickComponent] {
	return t.mouseDoubleRightClick
}
func (t *object) MouseEnter() ecs.EntityComponent[inputs.MouseEnterComponent] { return t.mouseEnter }
func (t *object) MouseLeave() ecs.EntityComponent[inputs.MouseLeaveComponent] { return t.mouseLeave }
func (t *object) MouseHover() ecs.EntityComponent[inputs.MouseHoverComponent] { return t.mouseHover }
func (t *object) MouseDrag() ecs.EntityComponent[inputs.MouseDragComponent]   { return t.mouseDrag }
