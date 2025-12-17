package mobilecamerasys

import (
	"engine/modules/camera"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type scrollSystem struct {
	window      window.Api
	cameraCtors camera.Interface
	logger      logger.Logger

	world             ecs.World
	transformTool     transform.Interface
	orthoArray        ecs.ComponentsArray[camera.OrthoComponent]
	mobileCameraArray ecs.ComponentsArray[camera.MobileCameraComponent]

	minZoom, maxZoom float32
}

func NewScrollSystem(
	logger logger.Logger,
	cameraCtors ecs.ToolFactory[camera.CameraTool],
	transformTool ecs.ToolFactory[transform.TransformTool],
	window window.Api,
	minZoom, maxZoom float32,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &scrollSystem{
			window:      window,
			cameraCtors: cameraCtors.Build(w).Camera(),
			logger:      logger,

			world:             w,
			transformTool:     transformTool.Build(w).Transform(),
			orthoArray:        ecs.GetComponentsArray[camera.OrthoComponent](w),
			mobileCameraArray: ecs.GetComponentsArray[camera.MobileCameraComponent](w),

			minZoom: minZoom, // e.g. 0.1
			maxZoom: maxZoom, // e.g. 5
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *scrollSystem) Listen(event sdl.MouseWheelEvent) {
	if event.Y == 0 {
		return
	}

	var mul = float32(math.Pow(10, float64(event.Y)/50))

	mousePos := s.window.GetMousePos()

	for _, cameraEntity := range s.mobileCameraArray.GetEntities() {
		ortho, ok := s.orthoArray.Get(cameraEntity)
		if !ok {
			continue
		}

		pos, _ := s.transformTool.AbsolutePos().Get(cameraEntity)
		rot, _ := s.transformTool.AbsoluteRotation().Get(cameraEntity)

		camera, err := s.cameraCtors.GetObject(cameraEntity)
		if err != nil {
			continue
		}

		rayBefore := camera.ShootRay(mousePos)

		// apply zoom
		ortho.Zoom *= mul
		ortho.Zoom = max(min(ortho.Zoom, s.maxZoom), s.minZoom)

		s.orthoArray.Set(cameraEntity, ortho)

		// read after
		rayAfter := camera.ShootRay(mousePos)

		// apply transform
		pos.Pos = pos.Pos.Add(rayBefore.Pos.Sub(rayAfter.Pos))

		rotationDifference := mgl32.QuatBetweenVectors(rayBefore.Direction, rayAfter.Direction)
		rot.Rotation = rotationDifference.Mul(rot.Rotation)

		s.transformTool.SetAbsolutePos(cameraEntity, pos)
		s.transformTool.SetAbsoluteRotation(cameraEntity, rot)
	}
}
