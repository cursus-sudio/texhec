package mouse

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"
	"errors"

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

	inputs.World
	inputs.InputsTool

	window window.Api

	maxMoved,
	moved float32 // max distance
	emitDrag     bool
	movingCamera ecs.EntityID
	movedEntity  *ecs.EntityID
	movedFrom    *window.MousePos
}

func NewClickSystem(
	logger logger.Logger,
	window window.Api,
	inputsToolFactory inputs.ToolFactory,
) inputs.System {
	return ecs.NewSystemRegister(func(w inputs.World) error {
		s := &clickSystem{
			logger:     logger,
			World:      w,
			InputsTool: inputsToolFactory.Build(w),

			window: window,

			maxMoved: 3,
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
		events.Emit(s.Events(), dragEvent)
	}

	if s.movedEntity != nil {
		entity := *s.movedEntity
		dragComponent, ok := s.Inputs().MouseDrag().Get(entity)
		if !ok {
			goto cleanUp
		}

		if i, ok := dragComponent.Event.(inputs.ApplyDragEvent); ok {
			dragComponent.Event = i.Apply(dragEvent)
		}
		events.EmitAny(s.Events(), dragComponent.Event)
	}

cleanUp:
	dist := mgl32.Vec2{
		float32(s.movedFrom.X - to.X),
		float32(s.movedFrom.Y - to.Y),
	}.Len()
	s.moved = dist + s.moved
	s.movedFrom = &to
}

func (s *clickSystem) ListenClick(event sdl.MouseButtonEvent) error {
	entities := s.Inputs().Hovered().GetEntities()
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
		s.moved = 0
		s.movedEntity = entity
		s.emitDrag = true

		if entity == nil {
			s.movedFrom = &pos
			break
		}

		if _, ok := s.Inputs().KeepSelected().Get(*entity); ok {
			s.emitDrag = false
			break
		}
		s.movedFrom = &pos
		hover, _ := s.Inputs().Hovered().Get(*entity)
		s.movingCamera = hover.Camera

	case sdl.RELEASED:
		dragged := s.movedEntity
		s.movedEntity = nil
		s.movedFrom = nil
		if entity == nil || dragged == nil || *entity != *dragged {
			break
		}

		if _, ok := s.Inputs().KeepSelected().Get(*entity); !ok && s.moved > s.maxMoved {
			break
		}

		var eventToEmit any

		switch event.Button {
		case sdl.BUTTON_LEFT:
			if comp, ok := s.Inputs().MouseLeft().Get(*entity); ok {
				eventToEmit = comp.Event
			}
			switch event.Clicks {
			case 2:
				if comp, ok := s.Inputs().MouseDoubleLeft().Get(*entity); ok {
					eventToEmit = comp.Event
				}
			}
		case sdl.BUTTON_RIGHT:
			if comp, ok := s.Inputs().MouseRight().Get(*entity); ok {
				eventToEmit = comp.Event
			}
			switch event.Clicks {
			case 2:
				if comp, ok := s.Inputs().MouseDoubleRight().Get(*entity); ok {
					eventToEmit = comp.Event
				}
			}
		}

		if eventToEmit != nil {
			events.EmitAny(s.Events(), eventToEmit)
		}
	}
	return nil
}
