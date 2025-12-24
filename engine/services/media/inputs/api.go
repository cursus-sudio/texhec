package inputs

import (
	"engine/services/clock"
	"engine/services/logger"
	"errors"
	"fmt"

	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	ErrNotHandledInput error = errors.New("not handled input")
)

type Api interface {
	Poll()
}

type api struct {
	logger logger.Logger
	clock  clock.Clock
	events events.Events
}

func newInputsApi(logger logger.Logger, clock clock.Clock, events events.Events) *api {
	return &api{
		logger: logger,
		clock:  clock,
		events: events,
	}
}

func (i *api) Poll() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		var e any
		switch event := event.(type) {
		case *sdl.AudioDeviceEvent:
			e = *event
		case *sdl.ClipboardEvent:
			e = *event
		case *sdl.CommonEvent:
			e = *event
		case *sdl.ControllerAxisEvent:
			e = *event
		case *sdl.ControllerButtonEvent:
			e = *event
		case *sdl.ControllerDeviceEvent:
			e = *event
		case *sdl.DisplayEvent:
			e = *event
		case *sdl.DollarGestureEvent:
			e = *event
		case *sdl.DropEvent:
			e = *event
		case *sdl.JoyAxisEvent:
			e = *event
		case *sdl.JoyBallEvent:
			e = *event
		case *sdl.JoyButtonEvent:
			e = *event
		case *sdl.JoyDeviceAddedEvent:
			e = *event
		case *sdl.JoyDeviceRemovedEvent:
			e = *event
		case *sdl.JoyHatEvent:
			e = *event
		case *sdl.KeyboardEvent:
			e = *event
		case *sdl.MouseButtonEvent:
			e = *event
		case *sdl.MouseMotionEvent:
			e = *event
		case *sdl.MouseWheelEvent:
			e = *event
		case *sdl.MultiGestureEvent:
			e = *event
		case *sdl.QuitEvent:
			e = *event
		case *sdl.RenderEvent:
			e = *event
		case *sdl.SensorEvent:
			e = *event
		case *sdl.TextInputEvent:
			e = *event
		case *sdl.TextEditingEvent:
			e = *event
		case *sdl.UserEvent:
			e = *event
		case *sdl.WindowEvent:
			e = *event
		case *sdl.TouchFingerEvent:
			e = *event
		default:
			i.logger.Warn(errors.Join(
				ErrNotHandledInput,
				fmt.Errorf("event not handled: type \"%d\": \"%v\"", event.GetType(), event),
			))
		}
		events.EmitAny(i.events, e)
	}
}
