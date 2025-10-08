package mousesys

import (
	"frontend/engine/components/mouse"
	"frontend/services/frames"
	"shared/services/ecs"

	"github.com/ogiusek/events"
)

type HoverEventSystem struct {
	world            ecs.World
	mouseEventsArray ecs.ComponentsArray[mouse.MouseEvents]
	events           events.Events
	query            ecs.LiveQuery
}

func NewHoverEventsSystem(world ecs.World, events events.Events) HoverEventSystem {
	query := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(mouse.MouseEvents{}),
		ecs.GetComponentType(mouse.Hovered{}),
	)
	return HoverEventSystem{
		world:            world,
		mouseEventsArray: ecs.GetComponentsArray[mouse.MouseEvents](world.Components()),
		events:           events,
		query:            query,
	}
}

func (s *HoverEventSystem) Listen(event frames.FrameEvent) {
	for _, entity := range s.query.Entities() {
		eventsComponent, err := s.mouseEventsArray.GetComponent(entity)
		if err != nil {
			continue
		}
		for _, event := range eventsComponent.HoverEvent {
			events.EmitAny(s.events, event)
		}
	}
}
