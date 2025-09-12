package mouse

import (
	"frontend/engine/components/mouse"
	"frontend/services/ecs"

	"github.com/ogiusek/events"
)

type HoverSystem struct {
	world            ecs.World
	mouseEventsArray ecs.ComponentsArray[mouse.MouseEvents]
	hoveredArray     ecs.ComponentsArray[mouse.Hovered]
	events           events.Events
	targets          map[ecs.ComponentType]ecs.EntityID
}

func NewHoverSystem(world ecs.World, events events.Events) HoverSystem {
	return HoverSystem{
		world:            world,
		mouseEventsArray: ecs.GetComponentsArray[mouse.MouseEvents](world.Components()),
		hoveredArray:     ecs.GetComponentsArray[mouse.Hovered](world.Components()),
		events:           events,
		targets:          map[ecs.ComponentType]ecs.EntityID{},
	}
}

func (s *HoverSystem) handleMouseLeave(entity ecs.EntityID) {
	s.hoveredArray.RemoveComponent(entity)

	mouseEvents, err := s.mouseEventsArray.GetComponent(entity)
	if err != nil {
		return
	}
	for _, event := range mouseEvents.MouseLeaveEvents {
		events.EmitAny(s.events, event)
	}
}

func (s *HoverSystem) Listen(event RayChangedTargetEvent) {
	if entity, ok := s.targets[event.ProjectionType]; ok {
		s.handleMouseLeave(entity)
	}
	if event.EntityID == nil {
		delete(s.targets, event.ProjectionType)
		return
	}
	s.targets[event.ProjectionType] = *event.EntityID
	entity := *event.EntityID

	mouseEvents, err := s.mouseEventsArray.GetComponent(entity)
	if err != nil {
		return
	}
	s.hoveredArray.SaveComponent(entity, mouse.Hovered{})
	for _, event := range mouseEvents.MouseEnterEvents {
		events.EmitAny(s.events, event)
	}
}
