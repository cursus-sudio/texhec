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
	clickOn uint8,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		liveQuery := w.QueryEntitiesWithComponents(
			ecs.GetComponentType(mouse.Hovered{}),
			ecs.GetComponentType(mouse.MouseEvents{}),
		)
		s := &clickSystem{
			world:            w,
			mouseEventsArray: ecs.GetComponentsArray[mouse.MouseEvents](w.Components()),
			events:           w.Events(),
			liveQuery:        liveQuery,
			clickOn:          clickOn,
		}

		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
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
