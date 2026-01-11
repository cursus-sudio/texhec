package mouse

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
)

type hoverEventSystem struct {
	world  ecs.World
	inputs inputs.Service

	events events.Events
}

func NewHoverEventsSystem(
	eventsBuilder events.Builder,
	world ecs.World,
	inputs inputs.Service,
) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := &hoverEventSystem{
			world:  world,
			inputs: inputs,
			events: eventsBuilder.Events(),
		}
		events.Listen(eventsBuilder, s.Listen)
		return nil
	})
}

func (s *hoverEventSystem) Listen(event frames.FrameEvent) {
	for _, entity := range s.inputs.Hovered().GetEntities() {
		eventsComponent, ok := s.inputs.Hover().Get(entity)
		if !ok {
			continue
		}

		events.EmitAny(s.events, eventsComponent.Event)
	}
}
