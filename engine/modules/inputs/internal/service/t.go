package service

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/events"
)

type RayChangedTargetEvent struct {
	Targets []inputs.Target
}

type service struct {
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
	eventsBuilder events.Builder,
	w ecs.World,
) inputs.Service {
	stack := []inputs.Target{}
	t := &service{
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
	events.Listen(eventsBuilder, func(e RayChangedTargetEvent) {
		*t.stackData = e.Targets
	})
	return t
}

func (t *service) Hovered() ecs.ComponentsArray[inputs.HoveredComponent] { return t.hovered }
func (t *service) Dragged() ecs.ComponentsArray[inputs.DraggedComponent] { return t.dragged }
func (t *service) Stacked() ecs.ComponentsArray[inputs.StackedComponent] { return t.stacked }

func (t *service) KeepSelected() ecs.ComponentsArray[inputs.KeepSelectedComponent] {
	return t.keepSelected
}

func (t *service) LeftClick() ecs.ComponentsArray[inputs.LeftClickComponent] { return t.mouseLeft }
func (t *service) DoubleLeftClick() ecs.ComponentsArray[inputs.DoubleLeftClickComponent] {
	return t.mouseDoubleLeft
}

func (t *service) RightClick() ecs.ComponentsArray[inputs.RightClickComponent] { return t.mouseRight }
func (t *service) DoubleRightClick() ecs.ComponentsArray[inputs.DoubleRightClickComponent] {
	return t.mouseDoubleRight
}

func (t *service) MouseEnter() ecs.ComponentsArray[inputs.MouseEnterComponent] { return t.mouseEnter }
func (t *service) MouseLeave() ecs.ComponentsArray[inputs.MouseLeaveComponent] { return t.mouseLeave }

func (t *service) Hover() ecs.ComponentsArray[inputs.HoverComponent] { return t.mouseHover }
func (t *service) Drag() ecs.ComponentsArray[inputs.DragComponent]   { return t.mouseDrag }

func (t *service) Stack() ecs.ComponentsArray[inputs.StackComponent] { return t.stack }

func (t *service) StackedData() []inputs.Target {
	stackCopy := make([]inputs.Target, len(*t.stackData))
	copy(stackCopy, *t.stackData)
	return stackCopy
}
