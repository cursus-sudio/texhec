package mobilecamerasys

import (
	"frontend/modules/camera"
	"frontend/modules/transform"
	"frontend/services/frames"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type wasdMoveSystem struct {
	logger        logger.Logger
	world         ecs.World
	transformTool transform.TransformTool
	orthoArray    ecs.ComponentsArray[camera.OrthoComponent]
	query         ecs.LiveQuery

	cameraCtors camera.CameraTool
	cameraSpeed float32
}

func NewWasdSystem(
	logger logger.Logger,
	cameraCtors ecs.ToolFactory[camera.CameraTool],
	transformToolFactory ecs.ToolFactory[transform.TransformTool],
	cameraSpeed float32,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		transformTool := transformToolFactory.Build(w)
		s := &wasdMoveSystem{
			logger:        logger,
			world:         w,
			transformTool: transformTool,
			orthoArray:    ecs.GetComponentsArray[camera.OrthoComponent](w.Components()),
			query: transformTool.Query(w.Query()).Require(
				ecs.GetComponentType(camera.OrthoComponent{}),
				ecs.GetComponentType(camera.MobileCameraComponent{}),
			).Build(),

			cameraCtors: cameraCtors.Build(w),
			cameraSpeed: cameraSpeed,
		}
		events.ListenE(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *wasdMoveSystem) Listen(event frames.FrameEvent) error {
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

	transformTransaction := s.transformTool.Transaction()

	for _, camera := range s.query.Entities() {
		transform := transformTransaction.GetEntity(camera)
		pos, err := transform.AbsolutePos().Get()
		if err != nil {
			s.logger.Warn(err)
			continue
		}

		ortho, err := s.orthoArray.GetComponent(camera)
		if err != nil {
			continue
		}

		pos.Pos = mgl32.Vec3{
			pos.Pos.X() + moveHorizontaly/ortho.Zoom,
			pos.Pos.Y() + moveVerticaly/ortho.Zoom,
			pos.Pos.Z(),
		}
		transform.AbsolutePos().Set(pos)
	}

	return ecs.FlushMany(transformTransaction.Transactions()...)
}
