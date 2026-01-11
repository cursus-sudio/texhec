package systems

import (
	"engine/modules/inputs"
	"engine/services/ecs"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type resizeSystem struct{}

func NewResizeSystem(eventsBuilder events.Builder) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := &resizeSystem{}
		events.Listen(eventsBuilder, s.Listen)
		return nil
	})
}

func (system *resizeSystem) Listen(e sdl.WindowEvent) {
	if e.Event != sdl.WINDOWEVENT_RESIZED {
		return
	}

	width, height := e.Data1, e.Data2
	gl.Viewport(0, 0, width, height)
}
