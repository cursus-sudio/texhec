package mouse

import (
	"engine/modules/inputs"
	"engine/modules/inputs/internal/tool"
	"engine/services/ecs"
	"engine/services/logger"
	"slices"

	"github.com/ogiusek/events"
)

type hoverSystem struct {
	inputs.World
	inputs.InputsTool
	logger  logger.Logger
	targets []inputs.Target
}

func NewHoverSystem(
	inputsToolFactory inputs.ToolFactory,
	logger logger.Logger,
) inputs.System {
	return ecs.NewSystemRegister(func(w inputs.World) error {
		s := &hoverSystem{
			World:      w,
			InputsTool: inputsToolFactory.Build(w),
			logger:     logger,
			targets:    nil,
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

func (s *hoverSystem) Listen(event tool.RayChangedTargetEvent) {
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
		s.Inputs().Hovered().Set(target.Entity, inputs.HoveredComponent{Camera: target.Camera})

		if mouseEnter, ok := s.Inputs().MouseEnter().Get(target.Entity); ok {
			events.EmitAny(s.Events(), mouseEnter.Event)
		}
	}
	s.targets = event.Targets

}
