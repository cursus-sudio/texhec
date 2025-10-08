package rendersys

import (
	"github.com/go-gl/gl/v4.5-core/gl"
)

type ClearSystem struct{}

func NewClearSystem() ClearSystem {
	return ClearSystem{}
}

func (s *ClearSystem) Listen(args RenderEvent) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
