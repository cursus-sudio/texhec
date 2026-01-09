package tool

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/logger"
	"sync"

	"github.com/ogiusek/events"
)

type RayChangedTargetEvent struct {
	Targets []inputs.Target
}

type tool struct {
	logger logger.Logger

	hovered ecs.ComponentsArray[inputs.HoveredComponent]
	dragged ecs.ComponentsArray[inputs.DraggedComponent]
	stacked ecs.ComponentsArray[inputs.StackedComponent]

	keepSelected ecs.ComponentsArray[inputs.KeepSelectedComponent]

	mouseLeft       ecs.ComponentsArray[inputs.LeftClickComponent]
	mouseDoubleLeft ecs.ComponentsArray[inputs.DoubleLeftClickComponent]

	mouseRight       ecs.ComponentsArray[inputs.RightClickComponent]
	mouseDoubleRight ecs.ComponentsArray[inputs.DoubleRightClickComponent]

	mouseEnter ecs.ComponentsArray[inputs.MouseEnterComponent]
	mouseLeave ecs.ComponentsArray[inputs.MouseLeaveComponent]

	mouseHover ecs.ComponentsArray[inputs.HoverComponent]
	mouseDrag  ecs.ComponentsArray[inputs.DragComponent]

	stack ecs.ComponentsArray[inputs.StackComponent]

	stackData *[]inputs.Target
}

func NewToolFactory(
	logger logger.Logger,
) inputs.ToolFactory {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w inputs.World) inputs.InputsTool {
		mutex.Lock()
		defer mutex.Unlock()
		if t, ok := ecs.GetGlobal[tool](w); ok {
			return t
		}
		stack := []inputs.Target{}
		t := &tool{
			logger,
			ecs.GetComponentsArray[inputs.HoveredComponent](w),
			ecs.GetComponentsArray[inputs.DraggedComponent](w),
			ecs.GetComponentsArray[inputs.StackedComponent](w),

			ecs.GetComponentsArray[inputs.KeepSelectedComponent](w),

			ecs.GetComponentsArray[inputs.LeftClickComponent](w),
			ecs.GetComponentsArray[inputs.DoubleLeftClickComponent](w),

			ecs.GetComponentsArray[inputs.RightClickComponent](w),
			ecs.GetComponentsArray[inputs.DoubleRightClickComponent](w),

			ecs.GetComponentsArray[inputs.MouseEnterComponent](w),
			ecs.GetComponentsArray[inputs.MouseLeaveComponent](w),

			ecs.GetComponentsArray[inputs.HoverComponent](w),
			ecs.GetComponentsArray[inputs.DragComponent](w),

			ecs.GetComponentsArray[inputs.StackComponent](w),

			&stack,
		}
		w.SaveGlobal(t)
		events.Listen(w.EventsBuilder(), func(e RayChangedTargetEvent) {
			*t.stackData = e.Targets
		})
		return t
	})
}

func (t *tool) Inputs() inputs.Interface {
	return t
}

func (t *tool) Hovered() ecs.ComponentsArray[inputs.HoveredComponent] { return t.hovered }
func (t *tool) Dragged() ecs.ComponentsArray[inputs.DraggedComponent] { return t.dragged }
func (t *tool) Stacked() ecs.ComponentsArray[inputs.StackedComponent] { return t.stacked }

func (t *tool) KeepSelected() ecs.ComponentsArray[inputs.KeepSelectedComponent] {
	return t.keepSelected
}

func (t *tool) LeftClick() ecs.ComponentsArray[inputs.LeftClickComponent] { return t.mouseLeft }
func (t *tool) DoubleLeftClick() ecs.ComponentsArray[inputs.DoubleLeftClickComponent] {
	return t.mouseDoubleLeft
}

func (t *tool) RightClick() ecs.ComponentsArray[inputs.RightClickComponent] { return t.mouseRight }
func (t *tool) DoubleRightClick() ecs.ComponentsArray[inputs.DoubleRightClickComponent] {
	return t.mouseDoubleRight
}

func (t *tool) MouseEnter() ecs.ComponentsArray[inputs.MouseEnterComponent] { return t.mouseEnter }
func (t *tool) MouseLeave() ecs.ComponentsArray[inputs.MouseLeaveComponent] { return t.mouseLeave }

func (t *tool) Hover() ecs.ComponentsArray[inputs.HoverComponent] { return t.mouseHover }
func (t *tool) Drag() ecs.ComponentsArray[inputs.DragComponent]   { return t.mouseDrag }

func (t *tool) Stack() ecs.ComponentsArray[inputs.StackComponent] { return t.stack }

func (t *tool) StackedData() []inputs.Target {
	stackCopy := make([]inputs.Target, len(*t.stackData))
	copy(stackCopy, *t.stackData)
	return stackCopy
}
