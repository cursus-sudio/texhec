package mobilecamerasys

import (
	"engine/modules/camera"
	"engine/modules/transform"
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

	world     ecs.World
	transform transform.Service
	camera    camera.Service

	minZoom, maxZoom float32
}

func NewScrollSystem(
	logger logger.Logger,
	window window.Api,
	eventsBuilder events.Builder,

	world ecs.World,
	transform transform.Service,
	camera camera.Service,

	minZoom, maxZoom float32,
) camera.System {
	return ecs.NewSystemRegister(func() error {
		s := &scrollSystem{
			window: window,
			logger: logger,

			world:     world,
			transform: transform,
			camera:    camera,

			minZoom: minZoom, // e.g. 0.1
			maxZoom: maxZoom, // e.g. 5
		}
		events.Listen(eventsBuilder, s.Listen)
		return nil
	})
}

func (s *scrollSystem) Listen(event sdl.MouseWheelEvent) {
	if event.Y == 0 {
		return
	}

	var mul = float32(math.Pow(10, float64(event.Y)/50))

	mousePos := s.window.GetMousePos()

	for _, cameraEntity := range s.camera.Mobile().GetEntities() {
		ortho, ok := s.camera.Ortho().Get(cameraEntity)
		if !ok {
			continue
		}

		pos, _ := s.transform.AbsolutePos().Get(cameraEntity)
		rot, _ := s.transform.AbsoluteRotation().Get(cameraEntity)

		rayBefore := s.camera.ShootRay(cameraEntity, mousePos)

		// apply zoom
		ortho.Zoom *= mul
		ortho.Zoom = max(min(ortho.Zoom, s.maxZoom), s.minZoom)

		s.camera.Ortho().Set(cameraEntity, ortho)

		// read after
		rayAfter := s.camera.ShootRay(cameraEntity, mousePos)

		// apply transform
		pos.Pos = pos.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		rot.Rotation = rotationDifference.Mul(rot.Rotation)

		s.transform.AbsolutePos().Set(cameraEntity, pos)
		s.transform.AbsoluteRotation().Set(cameraEntity, rot)
	}
}
