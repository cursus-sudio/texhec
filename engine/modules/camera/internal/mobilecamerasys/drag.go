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

	world         ecs.World
	transformTool transform.Tool
	query         ecs.LiveQuery

	cameraCtors camera.Tool
	window      window.Api
	logger      logger.Logger
}

func NewDragSystem(
	dragButton uint8,
	cameraCtors ecs.ToolFactory[camera.Tool],
	transformTool ecs.ToolFactory[transform.Tool],
	window window.Api,
	logger logger.Logger,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &dragSystem{
			isHeld: false,
			button: dragButton,

			world:         w,
			transformTool: transformTool.Build(w),
			query: w.Query().Require(
				ecs.GetComponentType(camera.MobileCameraComponent{}),
			).Build(),

			cameraCtors: cameraCtors.Build(w),
			window:      window,
			logger:      logger,
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *dragSystem) Listen(e inputs.DragEvent) {
	transformTransaction := s.transformTool.Transaction()

	for _, cameraEntity := range s.query.Entities() {
		transform := transformTransaction.GetObject(cameraEntity)
		pos, err := transform.AbsolutePos().Get()
		if err != nil {
			s.logger.Warn(err)
			continue
		}
		rot, err := transform.AbsoluteRotation().Get()
		if err != nil {
			s.logger.Warn(err)
			continue
		}

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

		transform.AbsolutePos().Set(pos)
		transform.AbsoluteRotation().Set(rot)
	}

	ecs.FlushMany(transformTransaction.Transactions()...)
}
