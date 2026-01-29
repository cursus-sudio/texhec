package systems

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/runtime"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type sys struct {
	Runtime       runtime.Runtime `inject:"1"`
	EventsBuilder events.Builder  `inject:"1"`
	Closed        bool
}

func NewQuitSystem(c ioc.Dic) inputs.ShutdownSystem {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*sys](c)
		events.Listen(s.EventsBuilder, s.Listen)
		events.Listen(s.EventsBuilder, s.ListenFrame)
		return nil
	})
}

func (s *sys) Listen(inputs.QuitEvent) {
	s.Closed = true
}

func (s *sys) ListenFrame(frames.FrameEvent) {
	if s.Closed {
		s.Runtime.Stop()
	}
}
