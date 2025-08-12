package inputs

import (
	"shared/services/runtime"

	"github.com/veandco/go-sdl2/sdl"
)

type QuitSystem struct {
	runtime runtime.Runtime
}

func NewQuitSystem(
	runtime runtime.Runtime,
) QuitSystem {
	return QuitSystem{
		runtime: runtime,
	}
}

func (system *QuitSystem) Listen(args sdl.QuitEvent) {
	system.runtime.Stop()
}
