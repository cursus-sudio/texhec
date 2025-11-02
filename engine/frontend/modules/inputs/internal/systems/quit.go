package systems

import (
	"frontend/modules/inputs"
	"shared/services/ecs"
	"shared/services/runtime"

	"github.com/ogiusek/events"
)

type sys struct {
	runtime runtime.Runtime
}

func NewQuitSystem(
	runtime runtime.Runtime,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &sys{
			runtime: runtime,
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *sys) Listen(inputs.QuitEvent) {
	s.runtime.Stop()
}
