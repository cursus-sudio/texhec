package render

import (
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/frames"
	"frontend/services/media/window"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

type RenderEvent struct{}

type RenderSystem struct {
	world  ecs.World
	events events.Events
	assets assets.Assets
	window window.Api

	setFence bool
	fence    uintptr
}

func NewRenderSystem(
	world ecs.World,
	events events.Events,
	window window.Api,
) RenderSystem {
	return RenderSystem{
		world:  world,
		events: events,
		window: window,
	}
}

func (s *RenderSystem) Listen(args frames.FrameEvent) error {
	if s.setFence {
		gl.WaitSync(s.fence, gl.SYNC_FLUSH_COMMANDS_BIT, gl.TIMEOUT_IGNORED)
		gl.DeleteSync(s.fence)
		s.setFence = false
		gl.Flush()
	}

	events.Emit(s.events, RenderEvent{})

	s.fence = gl.FenceSync(gl.SYNC_GPU_COMMANDS_COMPLETE, 0)
	s.setFence = true

	s.window.Window().GLSwap()

	return nil
}
