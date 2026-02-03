package systems

import (
	"engine/modules/camera"
	"engine/modules/render"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	"engine/services/media/window"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type renderSystem struct {
	World         ecs.World      `inject:"1"`
	Events        events.Events  `inject:"1"`
	Window        window.Api     `inject:"1"`
	EventsBuilder events.Builder `inject:"1"`
	Camera        camera.Service `inject:"1"`
	Logger        logger.Logger  `inject:"1"`
}

func NewRenderSystem(c ioc.Dic) render.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*renderSystem](c)
		events.ListenE(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *renderSystem) Listen(args frames.FrameEvent) error {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	cameras := s.Camera.OrderedCameras()
	for _, camera := range cameras {

		gl.Clear(gl.DEPTH_BUFFER_BIT)
		gl.Viewport(s.Camera.GetViewport(camera))

		events.Emit(s.Events, render.RenderEvent{
			Camera: camera,
		})
	}

	s.Window.Window().GLSwap()

	return nil
}
