package inputs

import (
	"shared/services/clock"
	"time"

	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

//

type EventType uint32

// events

// type OnCursorMoveEvent struct{ From, To CursorPos }
// type OnScrollEvent struct{ XMove, YMove int }
//
// type OnKeyPressEvent struct{ Key }
// type OnKeyReleaseEvent struct{ Key }
//
// type OnAxisMove struct {
// 	Axis     Axis
// 	From, To AxisState
// }

type InputsApi interface {
	IsHeld(Key) (duration time.Duration, isHeld bool)

	// axis // TODO when any game will need joystick
	// AxisPosition(Axis) AxisState

	// mouse
	CursorPos() CursorPos
}

type inputsApi struct {
	clock         clock.Clock
	events        events.Events
	held          map[Key]time.Time
	eventHandlers map[EventType]func(sdl.Event)
	cursorPos     CursorPos
}

func newInputsApi(clock clock.Clock, events events.Events) *inputsApi {
	x, y, _ := sdl.GetMouseState()

	return &inputsApi{
		clock:     clock,
		events:    events,
		held:      map[Key]time.Time{},
		cursorPos: CursorPos{X: int(x), Y: int(y)},
	}
}

func (i *inputsApi) Poll() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		eventType := EventType(event.GetType())
		eventHandler, ok := i.eventHandlers[eventType]
		if !ok {
			continue
		}
		eventHandler(event)
	}
}

func (i *inputsApi) IsHeld(key Key) (duration time.Duration, isHeld bool) {
	holdTime, ok := i.held[key]
	if !ok {
		return time.Duration(0), false
	}
	now := i.clock.Now()
	return now.Sub(holdTime), true
}

func (i *inputsApi) CursorPos() CursorPos {
	return i.cursorPos
}
