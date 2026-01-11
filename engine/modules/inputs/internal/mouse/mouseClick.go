package mouse

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"
	"errors"
	"slices"

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

	world  ecs.World
	inputs inputs.Service

	events events.Events
	window window.Api

	maxMoved,

	moved float32 // max distance

	emitDrag     bool
	movingCamera ecs.EntityID
	movedEntity  *ecs.EntityID
	movedFrom    *window.MousePos

	stacked []ecs.EntityID
}

func NewClickSystem(
	logger logger.Logger,
	window window.Api,
	eventsBuilder events.Builder,
	world ecs.World,
	inputs inputs.Service,
) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := &clickSystem{
			logger: logger,
			world:  world,
			inputs: inputs,

			events: eventsBuilder.Events(),
			window: window,

			maxMoved: 3,
		}

		events.Listen(eventsBuilder, s.ListenClick)
		events.Listen(eventsBuilder, s.ListenMove)
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
		events.Emit(s.events, dragEvent)
	}

	if s.movedEntity != nil {
		entity := *s.movedEntity
		dragComponent, ok := s.inputs.Drag().Get(entity)
		if !ok {
			goto cleanUp
		}

		if i, ok := dragComponent.Event.(inputs.ApplyDragEvent); ok {
			dragComponent.Event = i.Apply(dragEvent)
		}
		events.EmitAny(s.events, dragComponent.Event)
	}

cleanUp:
	dist := mgl32.Vec2{
		float32(s.movedFrom.X - to.X),
		float32(s.movedFrom.Y - to.Y),
	}.Len()
	s.moved = dist + s.moved
	s.movedFrom = &to
}

func (s *clickSystem) ListenClick(event sdl.MouseButtonEvent) {
	stackedBefore := make([]ecs.EntityID, len(s.stacked))
	copy(stackedBefore, s.stacked)

	stacked := []ecs.EntityID{}
	{
		for _, collision := range s.inputs.StackedData() {
			stacked = append(stacked, collision.Entity)
		}
	}

	var entity *ecs.EntityID

	i := 0
	for i = range s.stacked {
		if len(stacked) == i || stacked[i] != s.stacked[i] {
			break
		}
	}
	if len(s.stacked) != i && len(stacked) != i && stacked[i] == s.stacked[i] {
		i++
	}

	if i >= 0 && len(s.stacked) >= i && len(stacked) > i {
		s.stacked = s.stacked[:i]
		entity = &stacked[i]
	} else if len(stacked) != 0 {
		s.stacked = nil
		entity = &stacked[0]
	} else {
		s.stacked = nil
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

		if _, ok := s.inputs.KeepSelected().Get(*entity); ok {
			s.emitDrag = false
			break
		}
		s.movedFrom = &pos
		hover, _ := s.inputs.Hovered().Get(*entity)
		s.movingCamera = hover.Camera

	case sdl.RELEASED:
		dragged := s.movedEntity
		s.movedEntity = nil
		s.movedFrom = nil
		if entity == nil || dragged == nil || *entity != *dragged {
			break
		}

		if _, ok := s.inputs.KeepSelected().Get(*entity); !ok && s.moved > s.maxMoved {
			break
		}

		var eventToEmit any

		switch event.Button {
		case sdl.BUTTON_LEFT:
			if comp, ok := s.inputs.LeftClick().Get(*entity); ok {
				eventToEmit = comp.Event
			}
			switch event.Clicks {
			case 2:
				if comp, ok := s.inputs.DoubleLeftClick().Get(*entity); ok {
					eventToEmit = comp.Event
				}
			}
		case sdl.BUTTON_RIGHT:
			if comp, ok := s.inputs.RightClick().Get(*entity); ok {
				eventToEmit = comp.Event
			}
			switch event.Clicks {
			case 2:
				if comp, ok := s.inputs.DoubleRightClick().Get(*entity); ok {
					eventToEmit = comp.Event
				}
			}
		}

		if _, ok := s.inputs.Stack().Get(*entity); !ok {
			s.stacked = nil
		} else if len(s.stacked) != 0 && s.stacked[0] == *entity {
			s.stacked = s.stacked[:1]
		} else {
			s.stacked = append(s.stacked, *entity)
		}

		if eventToEmit != nil {
			events.EmitAny(s.events, eventToEmit)
		}
	}

	// find all added and removed
	removed := []ecs.EntityID{}
	for _, prevTarget := range stackedBefore {
		if slices.Contains(s.stacked, prevTarget) {
			continue
		}
		removed = append(removed, prevTarget)
	}
	added := []ecs.EntityID{}
	for _, target := range s.stacked {
		if slices.Contains(stackedBefore, target) {
			continue
		}
		added = append(added, target)
	}

	for _, added := range added {
		s.inputs.Stacked().Set(added, inputs.StackedComponent{})
	}

	for _, removed := range removed {
		s.inputs.Stacked().Remove(removed)
	}
}
