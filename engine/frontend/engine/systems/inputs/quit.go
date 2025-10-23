package inputssys

import (
	"shared/services/ecs"
	"shared/services/runtime"

	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type quitSystem struct {
	runtime runtime.Runtime
}

func NewQuitSystem(
	runtime runtime.Runtime,
) ecs.SystemRegister {
	return &quitSystem{
		runtime: runtime,
	}
}

func (s *quitSystem) Register(b events.Builder) {
	events.Listen(b, s.Listen)
}

func (system *quitSystem) Listen(args sdl.QuitEvent) {
	system.runtime.Stop()
}
