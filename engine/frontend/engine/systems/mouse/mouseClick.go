package mousesys

import (
	"frontend/engine/components/mouse"
	"shared/services/ecs"

	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type clickSystem struct {
	world            ecs.World
	mouseEventsArray ecs.ComponentsArray[mouse.MouseEvents]
	events           events.Events
	liveQuery        ecs.LiveQuery
	clickOn          uint8
}

// use sdl.PRESSED or sdl.RELEASED or clickOn
func NewClickSystem(
	world ecs.World,
	events events.Events,
	clickOn uint8,
) ecs.SystemRegister {
	liveQuery := world.QueryEntitiesWithComponents(
		ecs.GetComponentType(mouse.Hovered{}),
		ecs.GetComponentType(mouse.MouseEvents{}),
	)
	return &clickSystem{
		world:            world,
		mouseEventsArray: ecs.GetComponentsArray[mouse.MouseEvents](world.Components()),
		events:           events,
		liveQuery:        liveQuery,
		clickOn:          clickOn,
	}
}

func (s *clickSystem) Register(b events.Builder) {
	events.Listen(b, s.Listen)
}

func (s *clickSystem) Listen(event sdl.MouseButtonEvent) {
	if event.State != s.clickOn {
		return
	}
	entities := s.liveQuery.Entities()
	for _, entity := range entities {
		var mouseEvents mouse.MouseEvents
		mouseEvents, err := s.mouseEventsArray.GetComponent(entity)
		if err != nil {
			continue
		}

		var eventsToEmit []any

		switch event.Button {
		case sdl.BUTTON_LEFT:
			eventsToEmit = mouseEvents.LeftClickEvents
			switch event.Clicks {
			case 2:
				eventsToEmit = mouseEvents.DoubleLeftClickEvents
				break
			}
			break
		case sdl.BUTTON_RIGHT:
			eventsToEmit = mouseEvents.RightClickEvents
			switch event.Clicks {
			case 2:
				eventsToEmit = mouseEvents.DoubleRightClickEvents
				break
			}
			break
		}

		for _, event := range eventsToEmit {
			events.EmitAny(s.events, event)
		}
	}
}
