package mobilecamerasys

import (
	"frontend/engine/components/mobilecamera"
	"frontend/engine/components/transform"
	"frontend/engine/tools/cameras"
	"frontend/services/media/window"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type dragSystem struct {
	isHeld  bool
	button  uint8
	prevPos mgl32.Vec2

	world          ecs.World
	transformArray ecs.ComponentsArray[transform.Transform]
	query          ecs.LiveQuery

	cameraCtors cameras.CameraConstructors
	window      window.Api
	logger      logger.Logger
}

func NewDragSystem(
	dragButton uint8,
	cameraCtors cameras.CameraConstructorsFactory,
	window window.Api,
	logger logger.Logger,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &dragSystem{
			isHeld:  false,
			button:  dragButton,
			prevPos: mgl32.Vec2{0, 0},

			world:          w,
			transformArray: ecs.GetComponentsArray[transform.Transform](w.Components()),
			query: w.QueryEntitiesWithComponents(
				ecs.GetComponentType(mobilecamera.Component{}),
			),

			cameraCtors: cameraCtors.Build(w),
			window:      window,
			logger:      logger,
		}
		events.Listen(w.EventsBuilder(), s.Listen1)
		events.Listen(w.EventsBuilder(), s.Listen2)
		return nil
	})
}

func (s *dragSystem) Listen1(sdl.MouseMotionEvent) {
	prevPos := s.prevPos
	newPos := s.window.NormalizeMousePos(s.window.GetMousePos())
	s.prevPos = newPos

	if !s.isHeld {
		return
	}

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

func (s *dragSystem) Listen2(event sdl.MouseButtonEvent) {
	if event.Button != s.button {
		return
	}

	s.isHeld = event.State == sdl.PRESSED
}
