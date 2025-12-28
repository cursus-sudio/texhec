package mobilecamerasys

import (
	"engine/modules/camera"
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type dragSystem struct {
	isHeld bool
	button uint8

	camera.World
	camera.CameraTool

	window window.Api
	logger logger.Logger
}

func NewDragSystem(
	dragButton uint8,
	cameraCtors camera.ToolFactory,
	window window.Api,
	logger logger.Logger,
) camera.System {
	return ecs.NewSystemRegister(func(w camera.World) error {
		s := &dragSystem{
			isHeld: false,
			button: dragButton,

			World:      w,
			CameraTool: cameraCtors.Build(w),

			window: window,
			logger: logger,
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *dragSystem) Listen(e inputs.DragEvent) {
	for _, cameraEntity := range s.Camera().Mobile().GetEntities() {
		pos, _ := s.Transform().AbsolutePos().Get(cameraEntity)
		rot, _ := s.Transform().AbsoluteRotation().Get(cameraEntity)

		rayBefore := s.Camera().ShootRay(cameraEntity, e.From)
		rayAfter := s.Camera().ShootRay(cameraEntity, e.To)

		// apply difference
		pos.Pos = pos.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		rot.Rotation = rotationDifference.Mul(rot.Rotation)

		s.Transform().AbsolutePos().Set(cameraEntity, pos)
		s.Transform().AbsoluteRotation().Set(cameraEntity, rot)
	}
}
