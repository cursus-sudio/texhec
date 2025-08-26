package mouse

import (
	"frontend/engine/components/mouse"
	"frontend/services/ecs"
	"frontend/services/frames"

	"github.com/ogiusek/events"
)

type HoverEventSystem struct {
	world  ecs.World
	events events.Events
	query  ecs.LiveQuery
}

func NewHoverEventsSystem(world ecs.World, events events.Events) HoverEventSystem {
	query := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(mouse.MouseEvents{}),
		ecs.GetComponentType(mouse.Hovered{}),
	)
	return HoverEventSystem{
		world:  world,
		events: events,
		query:  query,
	}
}

func (s *HoverEventSystem) Listen(event frames.FrameEvent) {
	for _, entity := range s.query.Entities() {
		eventsComponent, err := ecs.GetComponent[mouse.MouseEvents](s.world, entity)
		if err != nil {
			continue
		}
		for _, event := range eventsComponent.HoverEvent {
			events.EmitAny(s.events, event)
		}
	}
}
