package mobilecamerasys

import (
	"engine/modules/camera"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type scrollSystem struct {
	window window.Api
	logger logger.Logger

	camera.World
	camera.CameraTool

	minZoom, maxZoom float32
}

func NewScrollSystem(
	logger logger.Logger,
	cameraCtors camera.ToolFactory,
	window window.Api,
	minZoom, maxZoom float32,
) camera.System {
	return ecs.NewSystemRegister(func(w camera.World) error {
		s := &scrollSystem{
			window: window,
			logger: logger,

			World:      w,
			CameraTool: cameraCtors.Build(w),

			minZoom: minZoom, // e.g. 0.1
			maxZoom: maxZoom, // e.g. 5
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *scrollSystem) Listen(event sdl.MouseWheelEvent) {
	if event.Y == 0 {
		return
	}

	var mul = float32(math.Pow(10, float64(event.Y)/50))

	mousePos := s.window.GetMousePos()

	for _, cameraEntity := range s.Camera().Mobile().GetEntities() {
		ortho, ok := s.Camera().Ortho().Get(cameraEntity)
		if !ok {
			continue
		}

		pos, _ := s.Transform().AbsolutePos().Get(cameraEntity)
		rot, _ := s.Transform().AbsoluteRotation().Get(cameraEntity)

		rayBefore := s.Camera().ShootRay(cameraEntity, mousePos)

		// apply zoom
		ortho.Zoom *= mul
		ortho.Zoom = max(min(ortho.Zoom, s.maxZoom), s.minZoom)

		s.Camera().Ortho().Set(cameraEntity, ortho)

		// read after
		rayAfter := s.Camera().ShootRay(cameraEntity, mousePos)

		// apply transform
		pos.Pos = pos.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		rot.Rotation = rotationDifference.Mul(rot.Rotation)

		s.Transform().AbsolutePos().Set(cameraEntity, pos)
		s.Transform().AbsoluteRotation().Set(cameraEntity, rot)
	}
}
