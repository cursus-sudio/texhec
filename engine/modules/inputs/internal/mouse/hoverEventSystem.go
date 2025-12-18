package mouse

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
)

type hoverEventSystem struct {
	inputs.World
	hoverEventArray ecs.ComponentsArray[inputs.MouseHoverComponent]
	hoveredArray    ecs.ComponentsArray[inputs.HoveredComponent]
}

func NewHoverEventsSystem() ecs.SystemRegister[inputs.World] {
	return ecs.NewSystemRegister(func(w inputs.World) error {
		s := &hoverEventSystem{
			World:           w,
			hoverEventArray: ecs.GetComponentsArray[inputs.MouseHoverComponent](w),
			hoveredArray:    ecs.GetComponentsArray[inputs.HoveredComponent](w),
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *hoverEventSystem) Listen(event frames.FrameEvent) {
	for _, entity := range s.hoveredArray.GetEntities() {
		eventsComponent, ok := s.hoverEventArray.Get(entity)
		if !ok {
			continue
		}

		events.EmitAny(s.Events(), eventsComponent.Event)
	}
}
