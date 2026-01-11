package systems

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/runtime"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type sys struct {
	Runtime       runtime.Runtime `inject:"1"`
	EventsBuilder events.Builder  `inject:"1"`
}

func NewQuitSystem(c ioc.Dic) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*sys](c)
		events.Listen(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *sys) Listen(inputs.QuitEvent) {
	s.Runtime.Stop()
}
