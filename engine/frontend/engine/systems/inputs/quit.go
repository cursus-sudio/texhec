package inputs

import (
	"shared/services/runtime"

	"github.com/veandco/go-sdl2/sdl"
)

type QuitSystem struct {
	Runtime runtime.Runtime
}

func NewQuitSystem(
	runtime runtime.Runtime,
) QuitSystem {
	return QuitSystem{
		Runtime: runtime,
	}
}

func (system *QuitSystem) Listen(args sdl.QuitEvent) {
	system.Runtime.Stop()
}
