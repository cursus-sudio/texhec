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
) ecs.SystemRegister[inputs.World] {
	return ecs.NewSystemRegister(func(w inputs.World) error {
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
