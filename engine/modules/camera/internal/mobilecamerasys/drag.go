package mobilecamerasys

import (
	"engine/modules/camera"
	"engine/modules/inputs"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type dragSystem struct {
	isHeld bool
	button uint8

	world     ecs.World
	transform transform.Service
	camera    camera.Service

	window window.Api
	logger logger.Logger
}

func NewDragSystem(
	dragButton uint8,

	world ecs.World,
	transform transform.Service,
	camera camera.Service,

	eventsBuilder events.Builder,
	window window.Api,
	logger logger.Logger,
) camera.System {
	return ecs.NewSystemRegister(func() error {
		s := &dragSystem{
			isHeld: false,
			button: dragButton,

			world:     world,
			transform: transform,
			camera:    camera,

			window: window,
			logger: logger,
		}
		events.Listen(eventsBuilder, s.Listen)
		return nil
	})
}

func (s *dragSystem) Listen(e inputs.DragEvent) {
	for _, cameraEntity := range s.camera.Mobile().GetEntities() {
		pos, _ := s.transform.AbsolutePos().Get(cameraEntity)
		rot, _ := s.transform.AbsoluteRotation().Get(cameraEntity)

		rayBefore := s.camera.ShootRay(cameraEntity, e.From)
		rayAfter := s.camera.ShootRay(cameraEntity, e.To)

		// apply difference
		pos.Pos = pos.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		rot.Rotation = rotationDifference.Mul(rot.Rotation)

		s.transform.AbsolutePos().Set(cameraEntity, pos)
		s.transform.AbsoluteRotation().Set(cameraEntity, rot)
	}
}
