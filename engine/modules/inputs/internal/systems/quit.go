package systems

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/runtime"

	"github.com/ogiusek/events"
)

type sys struct {
	runtime runtime.Runtime
}

func NewQuitSystem(
	runtime runtime.Runtime,
	eventsBuilder events.Builder,
) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := &sys{
			runtime: runtime,
		}
		events.Listen(eventsBuilder, s.Listen)
		return nil
	})
}

func (s *sys) Listen(inputs.QuitEvent) {
	s.runtime.Stop()
}
