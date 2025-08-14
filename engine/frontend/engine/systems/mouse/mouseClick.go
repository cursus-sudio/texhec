package mouse

import (
	"frontend/engine/components/mouse"
	"frontend/services/ecs"

	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type ClickSystem struct {
	world   ecs.World
	events  events.Events
	clickOn uint8
}

// use sdl.PRESSED or sdl.RELEASED or clickOn
func NewClickSystem(
	world ecs.World,
	events events.Events,
	clickOn uint8,
) ClickSystem {
	return ClickSystem{
		world:   world,
		events:  events,
		clickOn: clickOn,
	}
}

func (s ClickSystem) Listen(event sdl.MouseButtonEvent) {
	if event.State != s.clickOn {
		return
	}
	entities := s.world.GetEntitiesWithComponents(
		ecs.GetComponentType(mouse.Hovered{}),
		ecs.GetComponentType(mouse.MouseEvents{}),
	)
	for _, entity := range entities {
		var mouseEvents mouse.MouseEvents
		mouseEvents, err := ecs.GetComponent[mouse.MouseEvents](s.world, entity)
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
