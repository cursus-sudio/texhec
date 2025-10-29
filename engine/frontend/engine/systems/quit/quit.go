package quitsys

import (
	"shared/services/ecs"
	"shared/services/runtime"

	"github.com/ogiusek/events"
)

type QuitEvent struct{}

func NewQuitEvent() QuitEvent { return QuitEvent{} }

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

func (s *sys) Listen(QuitEvent) {
	s.runtime.Stop()
}
