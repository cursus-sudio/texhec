package mobilecamerasys

import (
	"engine/modules/camera"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type wasdMoveSystem struct {
	logger            logger.Logger
	world             ecs.World
	transformTool     transform.Interface
	orthoArray        ecs.ComponentsArray[camera.OrthoComponent]
	mobileCameraArray ecs.ComponentsArray[camera.MobileCameraComponent]

	cameraCtors camera.Camera
	cameraSpeed float32
}

func NewWasdSystem(
	logger logger.Logger,
	cameraCtors ecs.ToolFactory[camera.Camera],
	transformToolFactory ecs.ToolFactory[transform.Transform],
	cameraSpeed float32,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		transformTool := transformToolFactory.Build(w)
		s := &wasdMoveSystem{
			logger:            logger,
			world:             w,
			transformTool:     transformTool.Transform(),
			orthoArray:        ecs.GetComponentsArray[camera.OrthoComponent](w),
			mobileCameraArray: ecs.GetComponentsArray[camera.MobileCameraComponent](w),

			cameraCtors: cameraCtors.Build(w),
			cameraSpeed: cameraSpeed,
		}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *wasdMoveSystem) Listen(event frames.FrameEvent) {
	var moveVerticaly float32 = 0
	var moveHorizontaly float32 = 0
	{
		keys := sdl.GetKeyboardState()
		if keys[sdl.SCANCODE_W] == 1 || keys[sdl.SCANCODE_UP] == 1 {
			moveVerticaly = 1
		}
		if keys[sdl.SCANCODE_S] == 1 || keys[sdl.SCANCODE_DOWN] == 1 {
			moveVerticaly = -1
		}

		if keys[sdl.SCANCODE_A] == 1 || keys[sdl.SCANCODE_LEFT] == 1 {
			moveHorizontaly = -1
		}
		if keys[sdl.SCANCODE_D] == 1 || keys[sdl.SCANCODE_RIGHT] == 1 {
			moveHorizontaly = 1
		}
	}

	{
		moveHorizontaly *= float32(event.Delta.Milliseconds()) * s.cameraSpeed
		moveVerticaly *= float32(event.Delta.Milliseconds()) * s.cameraSpeed
	}

	for _, camera := range s.mobileCameraArray.GetEntities() {
		pos, _ := s.transformTool.AbsolutePos().GetComponent(camera)
		ortho, ok := s.orthoArray.GetComponent(camera)
		if !ok {
			continue
		}

		pos.Pos = mgl32.Vec3{
			pos.Pos.X() + moveHorizontaly/ortho.Zoom,
			pos.Pos.Y() + moveVerticaly/ortho.Zoom,
			pos.Pos.Z(),
		}
		s.transformTool.SetAbsolutePos(camera, pos)
	}
}
