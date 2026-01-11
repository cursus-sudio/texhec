package internal

import (
	"engine/modules/camera"
	"engine/modules/drag"
	"engine/modules/transform"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type s struct {
	eventsBuilder events.Builder
	transform     transform.Service
	camera        camera.Service
}

func NewSystem(
	eventsBuilder events.Builder,
	transform transform.Service,
	camera camera.Service,
) drag.System {
	return s{eventsBuilder, transform, camera}
}

func (s s) Register() error {
	events.Listen(s.eventsBuilder, func(event drag.DraggableEvent) {
		entity := event.Entity

		pos, _ := s.transform.AbsolutePos().Get(entity)
		rot, _ := s.transform.AbsoluteRotation().Get(entity)

		fromRay := s.camera.ShootRay(event.Drag.Camera, event.Drag.From)
		toRay := s.camera.ShootRay(event.Drag.Camera, event.Drag.To)

		posDiff := toRay.Pos.Sub(fromRay.Pos)
		pos.Pos = pos.Pos.Add(posDiff)
		s.transform.AbsolutePos().Set(entity, pos)

		rotDiff := mgl32.QuatBetweenVectors(toRay.Direction, fromRay.Direction)
		rot.Rotation = rot.Rotation.Mul(rotDiff)
		s.transform.AbsoluteRotation().Set(entity, rot)
	})
	return nil
}
