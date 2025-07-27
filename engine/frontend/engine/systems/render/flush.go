package render

import (
	"frontend/services/frames"
	"frontend/services/media/window"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type FlushSystem struct {
	Window window.Api
}

func NewFlushSystem(
	window window.Api,
) FlushSystem {
	return FlushSystem{
		Window: window,
	}
}

func (s *FlushSystem) Update(args frames.FrameEvent) {
	s.Window.Window().GLSwap()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
