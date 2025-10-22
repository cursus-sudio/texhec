package mobilecamerasys

import (
	"frontend/engine/components/mobilecamera"
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/engine/tools/cameras"
	"frontend/services/frames"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

type WasdMoveSystem struct {
	world                ecs.World
	transformArray       ecs.ComponentsArray[transform.Transform]
	orthoArray           ecs.ComponentsArray[projection.Ortho]
	transformTransaction ecs.ComponentsArrayTransaction[transform.Transform]
	query                ecs.LiveQuery

	cameraCtors cameras.CameraConstructors
	cameraSpeed float32
}

func NewWasdSystem(
	world ecs.World,
	cameraCtors cameras.CameraConstructors,
	cameraSpeed float32,
) WasdMoveSystem {
	return WasdMoveSystem{
		world:                world,
		transformArray:       ecs.GetComponentsArray[transform.Transform](world.Components()),
		orthoArray:           ecs.GetComponentsArray[projection.Ortho](world.Components()),
		transformTransaction: ecs.GetComponentsArray[transform.Transform](world.Components()).Transaction(),
		query: world.QueryEntitiesWithComponents(
			ecs.GetComponentType(transform.Transform{}),
			ecs.GetComponentType(projection.Ortho{}),
			ecs.GetComponentType(mobilecamera.Component{}),
		),

		cameraCtors: cameraCtors,
		cameraSpeed: cameraSpeed,
	}
}

func (s *WasdMoveSystem) Listen(event frames.FrameEvent) error {
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

	for _, camera := range s.query.Entities() {
		transform, err := s.transformArray.GetComponent(camera)
		if err != nil {
			continue
		}

		ortho, err := s.orthoArray.GetComponent(camera)
		if err != nil {
			continue
		}

		transform.SetPos(mgl32.Vec3{
			transform.Pos.X() + moveHorizontaly/ortho.Zoom,
			transform.Pos.Y() + moveVerticaly/ortho.Zoom,
			transform.Pos.Z(),
		})

		s.transformTransaction.SaveComponent(camera, transform)
	}

	return s.transformTransaction.Flush()
}
