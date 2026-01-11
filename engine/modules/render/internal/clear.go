package internal

import (
	"engine/modules/render"
	"engine/services/ecs"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

type clearSystem struct{}

func NewClearSystem(
	eventsBuilder events.Builder,
) render.System {
	return ecs.NewSystemRegister(func() error {
		s := &clearSystem{}
		events.Listen(eventsBuilder, s.Listen)
		return nil
	})
}

func (s *clearSystem) Listen(args render.RenderEvent) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
