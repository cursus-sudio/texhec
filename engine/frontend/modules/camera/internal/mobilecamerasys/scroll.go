package mobilecamerasys

import (
	"frontend/modules/camera"
	"frontend/modules/transform"
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
	cameraCtors camera.Tool
	logger      logger.Logger

	world         ecs.World
	transformTool transform.Tool
	orthoArray    ecs.ComponentsArray[camera.OrthoComponent]
	query         ecs.LiveQuery

	minZoom, maxZoom float32
}

func NewScrollSystem(
	logger logger.Logger,
	cameraCtors ecs.ToolFactory[camera.Tool],
	transformTool ecs.ToolFactory[transform.Tool],
	window window.Api,
	minZoom, maxZoom float32,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &scrollSystem{
			window:      window,
			cameraCtors: cameraCtors.Build(w),
			logger:      logger,

			world:         w,
			orthoArray:    ecs.GetComponentsArray[camera.OrthoComponent](w),
			transformTool: transformTool.Build(w),
			query: w.Query().Require(
				ecs.GetComponentType(camera.MobileCameraComponent{}),
			).Build(),

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

	mousePos := s.window.GetMousePos()

	transformTransaction := s.transformTool.Transaction()

	for _, cameraEntity := range s.query.Entities() {
		ortho, err := s.orthoArray.GetComponent(cameraEntity)
		if err != nil {
			continue
		}

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

		rayBefore := camera.ShootRay(mousePos)

		// apply zoom
		ortho.Zoom *= mul
		ortho.Zoom = max(min(ortho.Zoom, s.maxZoom), s.minZoom)

		if err := s.orthoArray.SaveComponent(cameraEntity, ortho); err != nil {
			return err
		}

		// read after
		rayAfter := camera.ShootRay(mousePos)

		// apply transform
		pos.Pos = pos.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		rot.Rotation = rotationDifference.Mul(rot.Rotation)

		transform.AbsolutePos().Set(pos)
		transform.AbsoluteRotation().Set(rot)
	}

	return ecs.FlushMany(transformTransaction.Transactions()...)
}
