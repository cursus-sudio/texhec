package internal

import (
	"engine/modules/camera"
	"engine/modules/drag"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type s struct {
	logger               logger.Logger
	cameraToolFactory    ecs.ToolFactory[camera.Tool]
	transformToolFactory ecs.ToolFactory[transform.Tool]
}

func NewSystem(
	logger logger.Logger,
	cameraToolFactory ecs.ToolFactory[camera.Tool],
	transformToolFactory ecs.ToolFactory[transform.Tool],
) ecs.SystemRegister {
	return s{logger, cameraToolFactory, transformToolFactory}
}

func (s s) Register(w ecs.World) error {
	cameraTool := s.cameraToolFactory.Build(w)
	transformTransaction := s.transformToolFactory.Build(w).Transaction()
	events.Listen(w.EventsBuilder(), func(event drag.DraggableEvent) {
		camera, err := cameraTool.GetObject(event.Drag.Camera)
		if err != nil {
			return
		}

		transform := transformTransaction.GetObject(event.Entity)
		pos, err := transform.AbsolutePos().Get()
		if err != nil {
			s.logger.Warn(err)
			return
		}
		rot, err := transform.AbsoluteRotation().Get()
		if err != nil {
			s.logger.Warn(err)
			return
		}

		fromRay := camera.ShootRay(event.Drag.From)
		toRay := camera.ShootRay(event.Drag.To)

		posDiff := toRay.Pos.Sub(fromRay.Pos)
		pos.Pos = pos.Pos.Add(posDiff)
		transform.AbsolutePos().Set(pos)

		rotDiff := mgl32.QuatBetweenVectors(toRay.Direction, fromRay.Direction)
		rot.Rotation = rot.Rotation.Mul(rotDiff)
		transform.AbsoluteRotation().Set(rot)
		s.logger.Warn(ecs.FlushMany(transformTransaction.Transactions()...))
	})
	return nil
}
