package mouse

import (
	"engine/modules/inputs"
	"engine/services/ecs"

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
	s.hoveredArray.Remove(entity)

	mouseLeave, ok := s.mouseLeaveArray.Get(entity)
	if !ok {
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

	s.hoveredArray.Set(entity, inputs.HoveredComponent{Camera: event.Camera})

	if mouseEnter, ok := s.mouseEnterArray.Get(entity); ok {
		events.EmitAny(s.events, mouseEnter.Event)
	}
}
