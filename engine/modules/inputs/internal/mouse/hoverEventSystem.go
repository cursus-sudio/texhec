package mouse

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
)

type hoverEventSystem struct {
	inputs.World
	inputs.InputsTool
}

func NewHoverEventsSystem(
	inputsToolFactory inputs.ToolFactory,
) inputs.System {
	return ecs.NewSystemRegister(func(w inputs.World) error {
		s := &hoverEventSystem{
			World:      w,
			InputsTool: inputsToolFactory.Build(w),
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *hoverEventSystem) Listen(event frames.FrameEvent) {
	for _, entity := range s.Inputs().Hovered().GetEntities() {
		eventsComponent, ok := s.Inputs().Hover().Get(entity)
		if !ok {
			continue
		}

		events.EmitAny(s.Events(), eventsComponent.Event)
	}
}
