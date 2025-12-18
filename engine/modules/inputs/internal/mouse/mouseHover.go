package mouse

import (
	"engine/modules/inputs"
	"engine/services/ecs"

	"github.com/ogiusek/events"
)

type hoverSystem struct {
	inputs.World
	inputs.InputsTool
	target *ecs.EntityID
}

func NewHoverSystem(
	inputsToolFactory ecs.ToolFactory[inputs.World, inputs.InputsTool],
) ecs.SystemRegister[inputs.World] {
	return ecs.NewSystemRegister(func(w inputs.World) error {
		s := &hoverSystem{
			World:      w,
			InputsTool: inputsToolFactory.Build(w),
			target:     nil,
		}

		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *hoverSystem) handleMouseLeave(entity ecs.EntityID) {
	s.Inputs().Hovered().Remove(entity)

	mouseLeave, ok := s.Inputs().MouseLeave().Get(entity)
	if !ok {
		return
	}
	events.EmitAny(s.Events(), mouseLeave.Event)
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

	s.Inputs().Hovered().Set(entity, inputs.HoveredComponent{Camera: event.Camera})

	if mouseEnter, ok := s.Inputs().MouseEnter().Get(entity); ok {
		events.EmitAny(s.Events(), mouseEnter.Event)
	}
}
