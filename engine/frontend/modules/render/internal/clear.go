package internal

import (
	"frontend/modules/render"
	"shared/services/ecs"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

type clearSystem struct{}

func NewClearSystem() ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &clearSystem{}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *clearSystem) Listen(args render.RenderEvent) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
