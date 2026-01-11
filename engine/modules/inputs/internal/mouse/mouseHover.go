package mouse

import (
	"engine/modules/inputs"
	"engine/modules/inputs/internal/service"
	"engine/services/ecs"
	"engine/services/logger"
	"slices"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type hoverSystem struct {
	EventsBuilder events.Builder `inject:"1"`
	Events        events.Events  `inject:"1"`
	World         ecs.World      `inject:"1"`
	Inputs        inputs.Service `inject:"1"`
	Logger        logger.Logger  `inject:"1"`
	targets       []inputs.Target
}

func NewHoverSystem(c ioc.Dic) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*hoverSystem](c)
		s.targets = nil

		events.Listen(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *hoverSystem) handleMouseLeave(entity ecs.EntityID) {
	s.Inputs.Hovered().Remove(entity)

	mouseLeave, ok := s.Inputs.MouseLeave().Get(entity)
	if !ok {
		return
	}
	events.EmitAny(s.Events, mouseLeave.Event)
}

func (s *hoverSystem) Listen(event service.RayChangedTargetEvent) {
	left := []inputs.Target{}
	entered := []inputs.Target{}

	for _, prevTarget := range s.targets {
		if slices.Contains(event.Targets, prevTarget) {
			continue
		}
		left = append(left, prevTarget)
	}
	for _, target := range event.Targets {
		if slices.Contains(s.targets, target) {
			continue
		}
		entered = append(entered, target)
	}

	for _, target := range left {
		s.handleMouseLeave(target.Entity)
	}

	for _, target := range entered {
		s.Inputs.Hovered().Set(target.Entity, inputs.HoveredComponent{Camera: target.Camera})

		if mouseEnter, ok := s.Inputs.MouseEnter().Get(target.Entity); ok {
			events.EmitAny(s.Events, mouseEnter.Event)
		}
	}
	s.targets = event.Targets

}
