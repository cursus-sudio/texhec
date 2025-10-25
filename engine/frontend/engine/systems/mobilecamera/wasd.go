package mobilecamerasys

import (
	"frontend/engine/components/mobilecamera"
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/engine/tools/cameras"
	"frontend/services/frames"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type wasdMoveSystem struct {
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
) ecs.SystemRegister {
	return &wasdMoveSystem{
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

func (s *wasdMoveSystem) Register(b events.Builder) {
	events.ListenE(b, s.Listen)
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

	for _, camera := range s.query.Entities() {
		transformComp, err := s.transformArray.GetComponent(camera)
		if err != nil {
			transformComp = transform.NewTransform()
		}

		ortho, err := s.orthoArray.GetComponent(camera)
		if err != nil {
			continue
		}

		transformComp.SetPos(mgl32.Vec3{
			transformComp.Pos.X() + moveHorizontaly/ortho.Zoom,
			transformComp.Pos.Y() + moveVerticaly/ortho.Zoom,
			transformComp.Pos.Z(),
		})

		s.transformTransaction.SaveComponent(camera, transformComp)
	}

	return s.transformTransaction.Flush()
}
