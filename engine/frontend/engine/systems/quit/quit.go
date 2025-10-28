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
	return &sys{
		runtime: runtime,
	}
}

func (s *sys) Register(b events.Builder) {
	events.Listen(b, s.Listen)
}

func (s *sys) Listen(QuitEvent) {
	s.runtime.Stop()
}
