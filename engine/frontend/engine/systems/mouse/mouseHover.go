package mouse

import (
	"frontend/engine/components/mouse"
	"frontend/services/ecs"

	"github.com/ogiusek/events"
)

type HoverSystem struct {
	world   ecs.World
	events  events.Events
	targets map[ecs.ComponentType]ecs.EntityID
}

func NewHoverSystem(world ecs.World, events events.Events) HoverSystem {
	return HoverSystem{
		world:   world,
		events:  events,
		targets: map[ecs.ComponentType]ecs.EntityID{},
	}
}

func (s *HoverSystem) handleMouseLeave(entity ecs.EntityID) {
	s.world.RemoveComponent(entity, ecs.GetComponentType(mouse.Hovered{}))

	mouseEvents, err := ecs.GetComponent[mouse.MouseEvents](s.world, entity)
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

	mouseEvents, err := ecs.GetComponent[mouse.MouseEvents](s.world, entity)
	if err != nil {
		return
	}
	s.world.SaveComponent(entity, mouse.Hovered{})
	for _, event := range mouseEvents.MouseEnterEvents {
		events.EmitAny(s.events, event)
	}
}
