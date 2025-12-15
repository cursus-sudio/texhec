package mouse

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
)

type hoverEventSystem struct {
	world           ecs.World
	hoverEventArray ecs.ComponentsArray[inputs.MouseHoverComponent]
	hoveredArray    ecs.ComponentsArray[inputs.HoveredComponent]
	events          events.Events
}

func NewHoverEventsSystem() ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &hoverEventSystem{
			world:           w,
			hoverEventArray: ecs.GetComponentsArray[inputs.MouseHoverComponent](w),
			hoveredArray:    ecs.GetComponentsArray[inputs.HoveredComponent](w),
			events:          w.Events(),
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *hoverEventSystem) Listen(event frames.FrameEvent) {
	for _, entity := range s.hoveredArray.GetEntities() {
		eventsComponent, ok := s.hoverEventArray.GetComponent(entity)
		if !ok {
			continue
		}

		events.EmitAny(s.events, eventsComponent.Event)
	}
}
