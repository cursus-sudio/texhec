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

	world             ecs.World
	transformTool     transform.Interface
	mobileCameraArray ecs.ComponentsArray[camera.MobileCameraComponent]

	cameraCtors camera.Interface
	window      window.Api
	logger      logger.Logger
}

func NewDragSystem(
	dragButton uint8,
	cameraCtors ecs.ToolFactory[camera.CameraTool],
	transformTool ecs.ToolFactory[transform.TransformTool],
	window window.Api,
	logger logger.Logger,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &dragSystem{
			isHeld: false,
			button: dragButton,

			world:             w,
			transformTool:     transformTool.Build(w).Transform(),
			mobileCameraArray: ecs.GetComponentsArray[camera.MobileCameraComponent](w),

			cameraCtors: cameraCtors.Build(w).Camera(),
			window:      window,
			logger:      logger,
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *dragSystem) Listen(e inputs.DragEvent) {
	for _, cameraEntity := range s.mobileCameraArray.GetEntities() {
		pos, _ := s.transformTool.AbsolutePos().Get(cameraEntity)
		rot, _ := s.transformTool.AbsoluteRotation().Get(cameraEntity)

		camera, err := s.cameraCtors.GetObject(cameraEntity)
		if err != nil {
			continue
		}
		rayBefore := camera.ShootRay(e.From)
		rayAfter := camera.ShootRay(e.To)

		// apply difference
		pos.Pos = pos.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		rot.Rotation = rotationDifference.Mul(rot.Rotation)

		s.transformTool.SetAbsolutePos(cameraEntity, pos)
		s.transformTool.SetAbsoluteRotation(cameraEntity, rot)
	}
}
