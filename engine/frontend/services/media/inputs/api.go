package inputs

import (
	"errors"
	"fmt"
	"shared/services/clock"
	"shared/services/logger"

	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	ErrNotHandledInput error = errors.New("not handled input")
)

type Api interface{}

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
		switch event.(type) {
		// TODO this creates an error
		case *sdl.AudioDeviceEvent:
			// events.Emit(i.events, *event.(*sdl.AudioDeviceEvent))
			e = *event.(*sdl.AudioDeviceEvent)
			break
		case *sdl.ClipboardEvent:
			e = *event.(*sdl.ClipboardEvent)
			break
		case *sdl.CommonEvent:
			e = *event.(*sdl.CommonEvent)
			break
		case *sdl.ControllerAxisEvent:
			e = *event.(*sdl.ControllerAxisEvent)
			break
		case *sdl.ControllerButtonEvent:
			e = *event.(*sdl.ControllerButtonEvent)
			break
		case *sdl.ControllerDeviceEvent:
			e = *event.(*sdl.ControllerDeviceEvent)
			break
		case *sdl.DisplayEvent:
			e = *event.(*sdl.DisplayEvent)
			break
		case *sdl.DollarGestureEvent:
			e = *event.(*sdl.DollarGestureEvent)
			break
		case *sdl.DropEvent:
			e = *event.(*sdl.DropEvent)
			break
		case *sdl.JoyAxisEvent:
			e = *event.(*sdl.JoyAxisEvent)
			break
		case *sdl.JoyBallEvent:
			e = *event.(*sdl.JoyBallEvent)
			break
		case *sdl.JoyButtonEvent:
			e = *event.(*sdl.JoyButtonEvent)
			break
		case *sdl.JoyDeviceAddedEvent:
			e = *event.(*sdl.JoyDeviceAddedEvent)
			break
		case *sdl.JoyDeviceRemovedEvent:
			e = *event.(*sdl.JoyDeviceRemovedEvent)
			break
		case *sdl.JoyHatEvent:
			e = *event.(*sdl.JoyHatEvent)
			break
		case *sdl.KeyboardEvent:
			e = *event.(*sdl.KeyboardEvent)
			break
		case *sdl.MouseButtonEvent:
			e = *event.(*sdl.MouseButtonEvent)
			break
		case *sdl.MouseMotionEvent:
			e = *event.(*sdl.MouseMotionEvent)
			break
		case *sdl.MouseWheelEvent:
			e = *event.(*sdl.MouseWheelEvent)
			break
		case *sdl.MultiGestureEvent:
			e = *event.(*sdl.MultiGestureEvent)
			break
		case *sdl.QuitEvent:
			e = *event.(*sdl.QuitEvent)
			break
		case *sdl.RenderEvent:
			e = *event.(*sdl.RenderEvent)
			break
		case *sdl.SensorEvent:
			e = *event.(*sdl.SensorEvent)
			break
		case *sdl.TextInputEvent:
			e = *event.(*sdl.TextInputEvent)
			break
		case *sdl.TextEditingEvent:
			e = *event.(*sdl.TextEditingEvent)
			break
		case *sdl.UserEvent:
			e = *event.(*sdl.UserEvent)
			break
		case *sdl.WindowEvent:
			e = *event.(*sdl.WindowEvent)
			break
		default:
			i.logger.Error(errors.Join(
				ErrNotHandledInput,
				fmt.Errorf("event not handled: type \"%d\": \"%v\"", event.GetType(), event),
			))
		}
		events.EmitAny(i.events, e)
	}
}
