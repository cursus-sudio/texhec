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

		camera, err := s.Camera().GetObject(cameraEntity)
		if err != nil {
			continue
		}
		rayBefore := camera.ShootRay(e.From)
		rayAfter := camera.ShootRay(e.To)

		// apply difference
		pos.Pos = pos.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		rot.Rotation = rotationDifference.Mul(rot.Rotation)

		s.Transform().SetAbsolutePos(cameraEntity, pos)
		s.Transform().SetAbsoluteRotation(cameraEntity, rot)
	}
}
