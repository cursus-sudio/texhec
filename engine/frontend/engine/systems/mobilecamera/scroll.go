package mobilecamerasys

import (
	"frontend/engine/components/mobilecamera"
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/engine/tools/cameras"
	"frontend/services/media/window"
	"math"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type scrollSystem struct {
	window      window.Api
	cameraCtors cameras.CameraResolver
	logger      logger.Logger

	world             ecs.World
	dynamicOrthoArray ecs.ComponentsArray[projection.DynamicOrtho]
	transformArray    ecs.ComponentsArray[transform.Transform]
	query             ecs.LiveQuery

	minZoom, maxZoom float32
}

func NewScrollSystem(
	logger logger.Logger,
	cameraCtors ecs.ToolFactory[cameras.CameraResolver],
	window window.Api,
	minZoom, maxZoom float32,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &scrollSystem{
			window:      window,
			cameraCtors: cameraCtors.Build(w),
			logger:      logger,

			world:             w,
			dynamicOrthoArray: ecs.GetComponentsArray[projection.DynamicOrtho](w.Components()),
			transformArray:    ecs.GetComponentsArray[transform.Transform](w.Components()),
			query: w.QueryEntitiesWithComponents(
				ecs.GetComponentType(mobilecamera.Component{}),
			),

			minZoom: minZoom, // e.g. 0.1
			maxZoom: maxZoom, // e.g. 5
		}
		events.ListenE(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *scrollSystem) Listen(event sdl.MouseWheelEvent) error {
	if event.Y == 0 {
		return nil
	}

	var mul = float32(math.Pow(10, float64(event.Y)/50))

	mousePos := s.window.NormalizeMousePos(s.window.GetMousePos())

	for _, cameraEntity := range s.query.Entities() {
		ortho, err := s.dynamicOrthoArray.GetComponent(cameraEntity)
		if err != nil {
			continue
		}

		transformComponent, err := s.transformArray.GetComponent(cameraEntity)
		if err != nil {
			transformComponent = transform.NewTransform()
		}

		camera, err := s.cameraCtors.Get(cameraEntity)
		if err != nil {
			continue
		}

		rayBefore := camera.ShootRay(mousePos)

		// apply zoom
		ortho.Zoom *= mul
		ortho.Zoom = max(min(ortho.Zoom, s.maxZoom), s.minZoom)

		if err := s.dynamicOrthoArray.SaveComponent(cameraEntity, ortho); err != nil {
			return err
		}

		// read after
		rayAfter := camera.ShootRay(mousePos)

		// apply transform
		transformComponent.Pos = transformComponent.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		transformComponent.Rotation = rotationDifference.Mul(transformComponent.Rotation)

		if err := s.transformArray.SaveComponent(cameraEntity, transformComponent); err != nil {
			return err
		}

	}

	return nil
}
