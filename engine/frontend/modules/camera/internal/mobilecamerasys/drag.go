package mobilecamerasys

import (
	"frontend/modules/camera"
	"frontend/modules/inputs"
	"frontend/modules/transform"
	"frontend/services/media/window"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type dragSystem struct {
	isHeld bool
	button uint8

	world          ecs.World
	transformArray ecs.ComponentsArray[transform.TransformComponent]
	query          ecs.LiveQuery

	cameraCtors camera.CameraTool
	window      window.Api
	logger      logger.Logger
}

func NewDragSystem(
	dragButton uint8,
	cameraCtors ecs.ToolFactory[camera.CameraTool],
	window window.Api,
	logger logger.Logger,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &dragSystem{
			isHeld: false,
			button: dragButton,

			world:          w,
			transformArray: ecs.GetComponentsArray[transform.TransformComponent](w.Components()),
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
	prevPos := s.window.NormalizeMousePos(int(e.From.X()), int(e.From.Y()))
	newPos := s.window.NormalizeMousePos(int(e.To.X()), int(e.To.Y()))

	for _, cameraEntity := range s.query.Entities() {
		transformComponent, err := s.transformArray.GetComponent(cameraEntity)
		if err != nil {
			transformComponent = transform.NewTransform()
		}

		camera, err := s.cameraCtors.Get(cameraEntity)
		if err != nil {
			continue
		}
		rayBefore := camera.ShootRay(prevPos)
		rayAfter := camera.ShootRay(newPos)

		// apply difference
		transformComponent.Pos = transformComponent.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		transformComponent.Rotation = rotationDifference.Mul(transformComponent.Rotation)

		s.transformArray.SaveComponent(cameraEntity, transformComponent)
	}
}
