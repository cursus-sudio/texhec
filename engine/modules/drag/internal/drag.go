package internal

import (
	"engine/modules/drag"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type s struct {
	logger logger.Logger
}

func NewSystem(
	logger logger.Logger,
) drag.System {
	return s{logger}
}

func (s s) Register(w drag.World) error {
	events.Listen(w.EventsBuilder(), func(event drag.DraggableEvent) {
		entity := event.Entity

		pos, _ := w.Transform().AbsolutePos().Get(entity)
		rot, _ := w.Transform().AbsoluteRotation().Get(entity)

		fromRay := w.Camera().ShootRay(event.Drag.Camera, event.Drag.From)
		toRay := w.Camera().ShootRay(event.Drag.Camera, event.Drag.To)

		posDiff := toRay.Pos.Sub(fromRay.Pos)
		pos.Pos = pos.Pos.Add(posDiff)
		w.Transform().AbsolutePos().Set(entity, pos)

		rotDiff := mgl32.QuatBetweenVectors(toRay.Direction, fromRay.Direction)
		rot.Rotation = rot.Rotation.Mul(rotDiff)
		w.Transform().AbsoluteRotation().Set(entity, rot)
	})
	return nil
}
