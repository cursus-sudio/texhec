package mouse

import (
	"engine/modules/inputs"
	"engine/modules/inputs/internal/service"
	"engine/services/ecs"
	"engine/services/logger"
	"slices"

	"github.com/ogiusek/events"
)

type hoverSystem struct {
	events  events.Events
	world   ecs.World
	inputs  inputs.Service
	logger  logger.Logger
	targets []inputs.Target
}

func NewHoverSystem(
	eventsBuilder events.Builder,
	world ecs.World,
	inputs inputs.Service,
	logger logger.Logger,
) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := &hoverSystem{
			events:  eventsBuilder.Events(),
			world:   world,
			inputs:  inputs,
			logger:  logger,
			targets: nil,
		}

		events.Listen(eventsBuilder, s.Listen)
		return nil
	})
}

func (s *hoverSystem) handleMouseLeave(entity ecs.EntityID) {
	s.inputs.Hovered().Remove(entity)

	mouseLeave, ok := s.inputs.MouseLeave().Get(entity)
	if !ok {
		return
	}
	events.EmitAny(s.events, mouseLeave.Event)
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
		s.inputs.Hovered().Set(target.Entity, inputs.HoveredComponent{Camera: target.Camera})

		if mouseEnter, ok := s.inputs.MouseEnter().Get(target.Entity); ok {
			events.EmitAny(s.events, mouseEnter.Event)
		}
	}
	s.targets = event.Targets

}
