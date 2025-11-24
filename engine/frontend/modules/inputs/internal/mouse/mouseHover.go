package mouse

import (
	"frontend/modules/inputs"
	"shared/services/ecs"

	"github.com/ogiusek/events"
)

type hoverSystem struct {
	world           ecs.World
	hoveredArray    ecs.ComponentsArray[inputs.HoveredComponent]
	mouseEnterArray ecs.ComponentsArray[inputs.MouseEnterComponent]
	mouseLeaveArray ecs.ComponentsArray[inputs.MouseLeaveComponent]
	events          events.Events
	target          *ecs.EntityID
}

func NewHoverSystem() ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &hoverSystem{
			world:           w,
			hoveredArray:    ecs.GetComponentsArray[inputs.HoveredComponent](w),
			mouseEnterArray: ecs.GetComponentsArray[inputs.MouseEnterComponent](w),
			mouseLeaveArray: ecs.GetComponentsArray[inputs.MouseLeaveComponent](w),
			events:          w.Events(),
			target:          nil,
		}

		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *hoverSystem) handleMouseLeave(entity ecs.EntityID) {
	s.hoveredArray.RemoveComponent(entity)

	mouseLeave, err := s.mouseLeaveArray.GetComponent(entity)
	if err != nil {
		return
	}
	events.EmitAny(s.events, mouseLeave.Event)
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

	s.hoveredArray.SaveComponent(entity, inputs.HoveredComponent{Camera: event.Camera})

	if mouseEnter, err := s.mouseEnterArray.GetComponent(entity); err == nil {
		events.EmitAny(s.events, mouseEnter.Event)
	}
}
