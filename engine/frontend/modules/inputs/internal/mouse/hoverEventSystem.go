package mouse

import (
	"frontend/modules/inputs"
	"frontend/services/frames"
	"shared/services/ecs"

	"github.com/ogiusek/events"
)

type hoverEventSystem struct {
	world           ecs.World
	hoverEventArray ecs.ComponentsArray[inputs.MouseHoverComponent]
	events          events.Events
	query           ecs.LiveQuery
}

func NewHoverEventsSystem() ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		query := w.Query().Require(
			ecs.GetComponentType(inputs.MouseHoverComponent{}),
			ecs.GetComponentType(inputs.HoveredComponent{}),
		).Build()
		s := &hoverEventSystem{
			world:           w,
			hoverEventArray: ecs.GetComponentsArray[inputs.MouseHoverComponent](w.Components()),
			events:          w.Events(),
			query:           query,
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *hoverEventSystem) Listen(event frames.FrameEvent) {
	for _, entity := range s.query.Entities() {
		eventsComponent, err := s.hoverEventArray.GetComponent(entity)
		if err != nil {
			continue
		}

		events.EmitAny(s.events, eventsComponent.Event)
	}
}
