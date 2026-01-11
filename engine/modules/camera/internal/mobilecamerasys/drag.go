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
	"github.com/ogiusek/ioc/v2"
)

type dragSystem struct {
	isHeld bool
	button uint8

	World     ecs.World         `inject:"1"`
	Transform transform.Service `inject:"1"`
	Camera    camera.Service    `inject:"1"`

	EventsBuilder events.Builder `inject:"1"`
	Window        window.Api     `inject:"1"`
	Logger        logger.Logger  `inject:"1"`
}

func NewDragSystem(c ioc.Dic, dragButton uint8) camera.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*dragSystem](c)
		s.isHeld = false
		s.button = dragButton
		events.Listen(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (s *dragSystem) Listen(e inputs.DragEvent) {
	for _, cameraEntity := range s.Camera.Mobile().GetEntities() {
		pos, _ := s.Transform.AbsolutePos().Get(cameraEntity)
		rot, _ := s.Transform.AbsoluteRotation().Get(cameraEntity)

		rayBefore := s.Camera.ShootRay(cameraEntity, e.From)
		rayAfter := s.Camera.ShootRay(cameraEntity, e.To)

		// apply difference
		pos.Pos = pos.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		rot.Rotation = rotationDifference.Mul(rot.Rotation)

		s.Transform.AbsolutePos().Set(cameraEntity, pos)
		s.Transform.AbsoluteRotation().Set(cameraEntity, rot)
	}
}
