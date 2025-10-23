package rendersys

import (
	"shared/services/ecs"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

type clearSystem struct{}

func NewClearSystem() ecs.SystemRegister {
	return &clearSystem{}
}

func (s *clearSystem) Register(b events.Builder) {
	events.Listen(b, s.Listen)
}

func (s *clearSystem) Listen(args RenderEvent) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
