package mouse

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"
	"errors"

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
	movedFrom    *window.MousePos
}

func NewClickSystem(logger logger.Logger, window window.Api) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &clickSystem{
			logger:       logger,
			world:        w,
			hoveredArray: ecs.GetComponentsArray[inputs.HoveredComponent](w),

			dragArray:             ecs.GetComponentsArray[inputs.MouseDragComponent](w),
			leftClickArray:        ecs.GetComponentsArray[inputs.MouseLeftClickComponent](w),
			doubleLeftClickArray:  ecs.GetComponentsArray[inputs.MouseDoubleLeftClickComponent](w),
			rightClickArray:       ecs.GetComponentsArray[inputs.MouseRightClickComponent](w),
			doubleRightClickArray: ecs.GetComponentsArray[inputs.MouseDoubleRightClickComponent](w),

			keepSelectedArray: ecs.GetComponentsArray[inputs.KeepSelectedComponent](w),

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
	to := window.NewMousePos(event.X, event.Y)
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
		dragComponent, ok := s.dragArray.Get(entity)
		if !ok {
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
	pos := window.NewMousePos(event.X, event.Y)

	switch event.State {
	case sdl.PRESSED:
		s.moved = false
		s.movedEntity = entity
		s.emitDrag = true

		if entity == nil {
			s.movedFrom = &pos
			break
		}

		if _, ok := s.keepSelectedArray.Get(*entity); ok {
			s.emitDrag = false
			break
		}
		s.movedFrom = &pos
		hover, _ := s.hoveredArray.Get(*entity)
		s.movingCamera = hover.Camera

	case sdl.RELEASED:
		dragged := s.movedEntity
		s.movedEntity = nil
		s.movedFrom = nil
		if entity == nil || dragged == nil || *entity != *dragged {
			break
		}

		if _, ok := s.keepSelectedArray.Get(*entity); !ok && s.moved {
			break
		}

		var eventToEmit any

		switch event.Button {
		case sdl.BUTTON_LEFT:
			if comp, ok := s.leftClickArray.Get(*entity); ok {
				eventToEmit = comp.Event
			}
			switch event.Clicks {
			case 2:
				if comp, ok := s.doubleLeftClickArray.Get(*entity); ok {
					eventToEmit = comp.Event
				}
				break
			}
			break
		case sdl.BUTTON_RIGHT:
			if comp, ok := s.rightClickArray.Get(*entity); ok {
				eventToEmit = comp.Event
			}
			switch event.Clicks {
			case 2:
				if comp, ok := s.doubleRightClickArray.Get(*entity); ok {
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
