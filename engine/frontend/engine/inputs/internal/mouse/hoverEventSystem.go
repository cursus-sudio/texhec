package mouse

import (
	"frontend/engine/inputs"
	"frontend/services/frames"
	"shared/services/ecs"

	"github.com/ogiusek/events"
)

type hoverEventSystem struct {
	world            ecs.World
	mouseEventsArray ecs.ComponentsArray[inputs.MouseEvents]
	events           events.Events
	query            ecs.LiveQuery
}

func NewHoverEventsSystem() ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		query := w.Query().Require(
			ecs.GetComponentType(inputs.MouseEvents{}),
			ecs.GetComponentType(inputs.Hovered{}),
		).Build()
		s := &hoverEventSystem{
			world:            w,
			mouseEventsArray: ecs.GetComponentsArray[inputs.MouseEvents](w.Components()),
			events:           w.Events(),
			query:            query,
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *hoverEventSystem) Listen(event frames.FrameEvent) {
	for _, entity := range s.query.Entities() {
		eventsComponent, err := s.mouseEventsArray.GetComponent(entity)
		if err != nil {
			continue
		}
		for _, event := range eventsComponent.HoverEvents {
			events.EmitAny(s.events, event)
		}
	}
}
