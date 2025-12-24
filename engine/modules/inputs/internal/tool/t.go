package tool

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"sync"
)

type tool struct {
	hovered ecs.ComponentsArray[inputs.HoveredComponent]
	dragged ecs.ComponentsArray[inputs.DraggedComponent]

	keepSelected ecs.ComponentsArray[inputs.KeepSelectedComponent]

	mouseLeft       ecs.ComponentsArray[inputs.MouseLeftClickComponent]
	mouseDoubleLeft ecs.ComponentsArray[inputs.MouseDoubleLeftClickComponent]

	mouseRight       ecs.ComponentsArray[inputs.MouseRightClickComponent]
	mouseDoubleRight ecs.ComponentsArray[inputs.MouseDoubleRightClickComponent]

	mouseEnter ecs.ComponentsArray[inputs.MouseEnterComponent]
	mouseLeave ecs.ComponentsArray[inputs.MouseLeaveComponent]

	mouseHover ecs.ComponentsArray[inputs.MouseHoverComponent]
	mouseDrag  ecs.ComponentsArray[inputs.MouseDragComponent]
}

func NewToolFactory() inputs.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w inputs.World) inputs.InputsTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		t := tool{
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
		w.SaveGlobal(t)
		return t
	})
}

func (t tool) Inputs() inputs.Interface {
	return t
}

func (t tool) Hovered() ecs.ComponentsArray[inputs.HoveredComponent] { return t.hovered }
func (t tool) Dragged() ecs.ComponentsArray[inputs.DraggedComponent] { return t.dragged }

func (t tool) KeepSelected() ecs.ComponentsArray[inputs.KeepSelectedComponent] { return t.keepSelected }

func (t tool) MouseLeft() ecs.ComponentsArray[inputs.MouseLeftClickComponent] { return t.mouseLeft }
func (t tool) MouseDoubleLeft() ecs.ComponentsArray[inputs.MouseDoubleLeftClickComponent] {
	return t.mouseDoubleLeft
}

func (t tool) MouseRight() ecs.ComponentsArray[inputs.MouseRightClickComponent] { return t.mouseRight }
func (t tool) MouseDoubleRight() ecs.ComponentsArray[inputs.MouseDoubleRightClickComponent] {
	return t.mouseDoubleRight
}

func (t tool) MouseEnter() ecs.ComponentsArray[inputs.MouseEnterComponent] { return t.mouseEnter }
func (t tool) MouseLeave() ecs.ComponentsArray[inputs.MouseLeaveComponent] { return t.mouseLeave }

func (t tool) MouseHover() ecs.ComponentsArray[inputs.MouseHoverComponent] { return t.mouseHover }
func (t tool) MouseDrag() ecs.ComponentsArray[inputs.MouseDragComponent]   { return t.mouseDrag }
