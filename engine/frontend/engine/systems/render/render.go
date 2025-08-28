package render

import (
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/frames"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

type RenderEvent struct{}

type RenderSystem struct {
	world  ecs.World
	events events.Events
	assets assets.Assets

	setFence bool
	fence    uintptr
}

func NewRenderSystem(
	world ecs.World,
	events events.Events,
) RenderSystem {
	return RenderSystem{
		world:  world,
		events: events,
	}
}

func (s *RenderSystem) Listen(args frames.FrameEvent) error {
	if s.setFence {
		s.setFence = false
		gl.ClientWaitSync(s.fence, gl.SYNC_FLUSH_COMMANDS_BIT, gl.TIMEOUT_IGNORED)
		gl.DeleteSync(s.fence)
	}

	events.Emit(s.events, RenderEvent{})

	s.fence = gl.FenceSync(gl.SYNC_GPU_COMMANDS_COMPLETE, 0)
	s.setFence = true

	return nil
}
