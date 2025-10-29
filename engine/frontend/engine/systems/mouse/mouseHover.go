package mousesys

import (
	"frontend/engine/components/mouse"
	"shared/services/ecs"

	"github.com/ogiusek/events"
)

type hoverSystem struct {
	world            ecs.World
	mouseEventsArray ecs.ComponentsArray[mouse.MouseEvents]
	hoveredArray     ecs.ComponentsArray[mouse.Hovered]
	events           events.Events
	target           *ecs.EntityID
}

func NewHoverSystem() ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &hoverSystem{
			world:            w,
			mouseEventsArray: ecs.GetComponentsArray[mouse.MouseEvents](w.Components()),
			hoveredArray:     ecs.GetComponentsArray[mouse.Hovered](w.Components()),
			events:           w.Events(),
			target:           nil,
		}

		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *hoverSystem) handleMouseLeave(entity ecs.EntityID) {
	s.hoveredArray.RemoveComponent(entity)

	mouseEvents, err := s.mouseEventsArray.GetComponent(entity)
	if err != nil {
		return
	}
	for _, event := range mouseEvents.MouseLeaveEvents {
		events.EmitAny(s.events, event)
	}
}

func (s *hoverSystem) Listen(event RayChangedTargetEvent) {
	if s.target != nil {
		s.handleMouseLeave(*s.target)
	}
	if event.EntityID == nil {
		s.target = nil
		return
	}
	s.target = event.EntityID
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
