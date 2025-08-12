package inputs

import (
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

type ResizeSystem struct{}

func NewResizeSystem() ResizeSystem {
	return ResizeSystem{}
}

func (system ResizeSystem) Listen(e sdl.WindowEvent) {
	if e.Event != sdl.WINDOWEVENT_RESIZED {
		return
	}

	width, height := e.Data1, e.Data2
	gl.Viewport(0, 0, width, height)
}
