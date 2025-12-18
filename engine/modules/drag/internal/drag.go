package internal

import (
	"engine/modules/drag"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type s struct {
	logger logger.Logger
}

func NewSystem(
	logger logger.Logger,
) ecs.SystemRegister[drag.World] {
	return s{logger}
}

func (s s) Register(w drag.World) error {
	events.Listen(w.EventsBuilder(), func(event drag.DraggableEvent) {
		camera, err := w.Camera().GetObject(event.Drag.Camera)
		if err != nil {
			return
		}
		entity := event.Entity

		pos, _ := w.Transform().AbsolutePos().Get(entity)
		rot, _ := w.Transform().AbsoluteRotation().Get(entity)

		fromRay := camera.ShootRay(event.Drag.From)
		toRay := camera.ShootRay(event.Drag.To)

		posDiff := toRay.Pos.Sub(fromRay.Pos)
		pos.Pos = pos.Pos.Add(posDiff)
		w.Transform().SetAbsolutePos(entity, pos)

		rotDiff := mgl32.QuatBetweenVectors(toRay.Direction, fromRay.Direction)
		rot.Rotation = rot.Rotation.Mul(rotDiff)
		w.Transform().SetAbsoluteRotation(entity, rot)
	})
	return nil
}
