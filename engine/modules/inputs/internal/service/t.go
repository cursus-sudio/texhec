package service

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type RayChangedTargetEvent struct {
	Targets []inputs.Target
}

type service struct {
	Logger        logger.Logger  `inject:"1"`
	World         ecs.World      `inject:"1"`
	EventsBuilder events.Builder `inject:"1"`

	hovered ecs.ComponentsArray[inputs.HoveredComponent]
	dragged ecs.ComponentsArray[inputs.DraggedComponent]
	stacked ecs.ComponentsArray[inputs.StackedComponent]

	keepSelected ecs.ComponentsArray[inputs.KeepSelectedComponent]

	leftClick       ecs.ComponentsArray[inputs.LeftClickComponent]
	doubleLeftClick ecs.ComponentsArray[inputs.DoubleLeftClickComponent]

	rightClick       ecs.ComponentsArray[inputs.RightClickComponent]
	doubleRightClick ecs.ComponentsArray[inputs.DoubleRightClickComponent]

	mouseEnter ecs.ComponentsArray[inputs.MouseEnterComponent]
	mouseLeave ecs.ComponentsArray[inputs.MouseLeaveComponent]

	mouseHover ecs.ComponentsArray[inputs.HoverComponent]
	mouseDrag  ecs.ComponentsArray[inputs.DragComponent]

	stack ecs.ComponentsArray[inputs.StackComponent]

	stackData *[]inputs.Target
}

func NewService(c ioc.Dic) inputs.Service {
	t := ioc.GetServices[*service](c)
	t.hovered = ecs.GetComponentsArray[inputs.HoveredComponent](t.World)
	t.dragged = ecs.GetComponentsArray[inputs.DraggedComponent](t.World)
	t.stacked = ecs.GetComponentsArray[inputs.StackedComponent](t.World)

	t.keepSelected = ecs.GetComponentsArray[inputs.KeepSelectedComponent](t.World)

	t.leftClick = ecs.GetComponentsArray[inputs.LeftClickComponent](t.World)
	t.doubleLeftClick = ecs.GetComponentsArray[inputs.DoubleLeftClickComponent](t.World)

	t.rightClick = ecs.GetComponentsArray[inputs.RightClickComponent](t.World)
	t.doubleRightClick = ecs.GetComponentsArray[inputs.DoubleRightClickComponent](t.World)

	t.mouseEnter = ecs.GetComponentsArray[inputs.MouseEnterComponent](t.World)
	t.mouseLeave = ecs.GetComponentsArray[inputs.MouseLeaveComponent](t.World)

	t.mouseHover = ecs.GetComponentsArray[inputs.HoverComponent](t.World)
	t.mouseDrag = ecs.GetComponentsArray[inputs.DragComponent](t.World)

	t.stack = ecs.GetComponentsArray[inputs.StackComponent](t.World)

	ecs.GetComponentsArray[inputs.StackComponent](t.World)

	stack := []inputs.Target{}
	t.stackData = &stack
	events.Listen(t.EventsBuilder, func(e RayChangedTargetEvent) {
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

func (t *service) LeftClick() ecs.ComponentsArray[inputs.LeftClickComponent] { return t.leftClick }
func (t *service) DoubleLeftClick() ecs.ComponentsArray[inputs.DoubleLeftClickComponent] {
	return t.doubleLeftClick
}

func (t *service) RightClick() ecs.ComponentsArray[inputs.RightClickComponent] { return t.rightClick }
func (t *service) DoubleRightClick() ecs.ComponentsArray[inputs.DoubleRightClickComponent] {
	return t.doubleRightClick
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
