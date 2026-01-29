package systems

import (
	"engine/modules/render"
	"engine/services/ecs"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type clearSystem struct {
	EventsBuilder events.Builder `inject:"1"`
}

func NewClearSystem(c ioc.Dic) render.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*clearSystem](c)
		events.Listen(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *clearSystem) Listen(args render.RenderEvent) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
