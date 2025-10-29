package rendersys

import (
	"frontend/services/assets"
	"frontend/services/frames"
	"frontend/services/media/window"
	"shared/services/ecs"
	"sync"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
)

type RenderEvent struct{}

type renderSystem struct {
	world  ecs.World
	events events.Events
	assets assets.Assets
	window window.Api

	fences       []uintptr
	buffersCount int
	mutex        sync.Locker
}

func NewRenderSystem(
	window window.Api,
	bufferCount int,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &renderSystem{
			world:  w,
			events: w.Events(),
			window: window,

			fences:       []uintptr{},
			buffersCount: max(1, bufferCount),
			mutex:        &sync.Mutex{},
		}
		events.ListenE(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *renderSystem) Listen(args frames.FrameEvent) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.fences) == s.buffersCount {
		fence := s.fences[0]
		s.fences = s.fences[1:]
		gl.WaitSync(fence, gl.SYNC_FLUSH_COMMANDS_BIT, gl.TIMEOUT_IGNORED)
		gl.DeleteSync(fence)
	}

	events.Emit(s.events, RenderEvent{})

	s.window.Window().GLSwap()

	fence := gl.FenceSync(gl.SYNC_GPU_COMMANDS_COMPLETE, 0)
	s.fences = append(s.fences, fence)

	return nil
}
