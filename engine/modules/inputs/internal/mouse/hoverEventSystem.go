package mouse

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type hoverEventSystem struct {
	World  ecs.World      `inject:"1"`
	Inputs inputs.Service `inject:"1"`

	EventsBuilder events.Builder `inject:"1"`
	Events        events.Events  `inject:"1"`
}

func NewHoverEventsSystem(c ioc.Dic) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*hoverEventSystem](c)
		events.Listen(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *hoverEventSystem) Listen(event frames.FrameEvent) {
	for _, entity := range s.Inputs.Hovered().GetEntities() {
		eventsComponent, ok := s.Inputs.Hover().Get(entity)
		if !ok {
			continue
		}

		events.EmitAny(s.Events, eventsComponent.Event)
	}
}
