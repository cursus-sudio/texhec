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
	cameraToolFactory    ecs.ToolFactory[camera.Camera]
	transformToolFactory ecs.ToolFactory[transform.Transform]
}

func NewSystem(
	logger logger.Logger,
	cameraToolFactory ecs.ToolFactory[camera.Camera],
	transformToolFactory ecs.ToolFactory[transform.Transform],
) ecs.SystemRegister {
	return s{logger, cameraToolFactory, transformToolFactory}
}

func (s s) Register(w ecs.World) error {
	cameraTool := s.cameraToolFactory.Build(w).Camera()
	transformTool := s.transformToolFactory.Build(w).Transform()
	events.Listen(w.EventsBuilder(), func(event drag.DraggableEvent) {
		camera, err := cameraTool.GetObject(event.Drag.Camera)
		if err != nil {
			return
		}
		entity := event.Entity

		pos, _ := transformTool.AbsolutePos().GetComponent(entity)
		rot, _ := transformTool.AbsoluteRotation().GetComponent(entity)

		fromRay := camera.ShootRay(event.Drag.From)
		toRay := camera.ShootRay(event.Drag.To)

		posDiff := toRay.Pos.Sub(fromRay.Pos)
		pos.Pos = pos.Pos.Add(posDiff)
		transformTool.SetAbsolutePos(entity, pos)

		rotDiff := mgl32.QuatBetweenVectors(toRay.Direction, fromRay.Direction)
		rot.Rotation = rot.Rotation.Mul(rotDiff)
		transformTool.SetAbsoluteRotation(entity, rot)
	})
	return nil
}
