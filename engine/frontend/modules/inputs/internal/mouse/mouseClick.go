package mouse

import (
	"errors"
	"frontend/modules/inputs"
	"frontend/services/media/window"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

//

var (
	ErrCanHoverOverMaxOneEntity error = errors.New("can hover over max one entity at a time")
)

type clickSystem struct {
	logger logger.Logger

	world        ecs.World
	hoveredArray ecs.ComponentsArray[inputs.HoveredComponent]

	dragArray             ecs.ComponentsArray[inputs.MouseDragComponent]
	leftClickArray        ecs.ComponentsArray[inputs.MouseLeftClickComponent]
	doubleLeftClickArray  ecs.ComponentsArray[inputs.MouseDoubleLeftClickComponent]
	rightClickArray       ecs.ComponentsArray[inputs.MouseRightClickComponent]
	doubleRightClickArray ecs.ComponentsArray[inputs.MouseDoubleRightClickComponent]

	keepSelectedArray ecs.ComponentsArray[inputs.KeepSelectedComponent]

	window window.Api

	moved        bool
	emitDrag     bool
	movingCamera ecs.EntityID
	movedEntity  *ecs.EntityID
	movedFrom    *mgl32.Vec2
}

func NewClickSystem(logger logger.Logger, window window.Api) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &clickSystem{
			logger:       logger,
			world:        w,
			hoveredArray: ecs.GetComponentsArray[inputs.HoveredComponent](w.Components()),

			dragArray:             ecs.GetComponentsArray[inputs.MouseDragComponent](w.Components()),
			leftClickArray:        ecs.GetComponentsArray[inputs.MouseLeftClickComponent](w.Components()),
			doubleLeftClickArray:  ecs.GetComponentsArray[inputs.MouseDoubleLeftClickComponent](w.Components()),
			rightClickArray:       ecs.GetComponentsArray[inputs.MouseRightClickComponent](w.Components()),
			doubleRightClickArray: ecs.GetComponentsArray[inputs.MouseDoubleRightClickComponent](w.Components()),

			keepSelectedArray: ecs.GetComponentsArray[inputs.KeepSelectedComponent](w.Components()),

			window: window,

			movedEntity: nil,
			movedFrom:   nil,
		}

		events.ListenE(w.EventsBuilder(), s.ListenClick)
		events.Listen(w.EventsBuilder(), s.ListenMove)
		return nil
	})
}

func (s *clickSystem) ListenMove(event sdl.MouseMotionEvent) {
	// s.draggedFrom has to be set to allow dragging
	if s.movedFrom == nil {
		return
	}

	from := *s.movedFrom
	to := s.window.NormalizeMousePos(int(event.X), int(event.Y))
	dragEvent := inputs.DragEvent{
		Camera: s.movingCamera,
		From:   from,
		To:     to,
	}

	if s.emitDrag {
		events.Emit(s.world.Events(), dragEvent)
	}

	if s.movedEntity != nil {
		entity := *s.movedEntity
		dragComponent, err := s.dragArray.GetComponent(entity)
		if err != nil {
			goto cleanUp
		}

		if i, ok := dragComponent.Event.(inputs.ApplyDragEvent); ok {
			dragComponent.Event = i.Apply(dragEvent)
		}
		events.EmitAny(s.world.Events(), dragComponent.Event)
	}

cleanUp:
	s.moved = true
	s.movedFrom = &to
}

func (s *clickSystem) ListenClick(event sdl.MouseButtonEvent) error {
	entities := s.hoveredArray.GetEntities()
	if len(entities) > 1 {
		return ErrCanHoverOverMaxOneEntity
	}

	var entity *ecs.EntityID
	if len(entities) == 1 {
		e := entities[0]
		entity = &e
	}
	pos := s.window.NormalizeMousePos(int(event.X), int(event.Y))

	switch event.State {
	case sdl.PRESSED:
		s.moved = false
		s.movedEntity = entity
		s.emitDrag = true

		if entity == nil {
			s.movedFrom = &pos
			break
		}

		if _, err := s.keepSelectedArray.GetComponent(*entity); err == nil {
			s.emitDrag = false
			break
		}
		s.movedFrom = &pos
		hover, _ := s.hoveredArray.GetComponent(*entity)
		s.movingCamera = hover.Camera

	case sdl.RELEASED:
		dragged := s.movedEntity
		s.movedEntity = nil
		s.movedFrom = nil
		if entity == nil || dragged == nil || *entity != *dragged {
			break
		}

		if _, err := s.keepSelectedArray.GetComponent(*entity); err != nil && s.moved {
			break
		}

		var eventToEmit any

		switch event.Button {
		case sdl.BUTTON_LEFT:
			if comp, err := s.leftClickArray.GetComponent(*entity); err == nil {
				eventToEmit = comp.Event
			}
			switch event.Clicks {
			case 2:
				if comp, err := s.doubleLeftClickArray.GetComponent(*entity); err == nil {
					eventToEmit = comp.Event
				}
				break
			}
			break
		case sdl.BUTTON_RIGHT:
			if comp, err := s.rightClickArray.GetComponent(*entity); err == nil {
				eventToEmit = comp.Event
			}
			switch event.Clicks {
			case 2:
				if comp, err := s.doubleRightClickArray.GetComponent(*entity); err == nil {
					eventToEmit = comp.Event
				}
				break
			}
			break
		}

		if eventToEmit != nil {
			events.EmitAny(s.world.Events(), eventToEmit)
		}
	}
	return nil
}
