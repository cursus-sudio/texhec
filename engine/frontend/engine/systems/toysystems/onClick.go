package toysystems

import (
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/services/colliders/shapes"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/graphics/camera"
	"frontend/services/media/window"
)

type ShootRayEvent struct{ X, Y int }

func NewShootRayEvent(x, y int) ShootRayEvent {
	return ShootRayEvent{X: x, Y: y}
}

type ShootRaySystem[Projection projection.Projection] struct {
	world          ecs.World
	window         window.Api
	console        console.Console
	distFromCamera float32
	getEntity      func() ecs.EntityId
}

func NewShootRaySystem[Projection projection.Projection](
	world ecs.World,
	window window.Api,
	console console.Console,
	distFromCamera float32,
	getEntity func() ecs.EntityId,
) ShootRaySystem[Projection] {
	return ShootRaySystem[Projection]{
		world:          world,
		window:         window,
		console:        console,
		distFromCamera: distFromCamera,
		getEntity:      func() ecs.EntityId { return getEntity() },
	}
}

func (system *ShootRaySystem[Projection]) Listen(args ShootRayEvent) error {
	var cameraTransform transform.Transform
	var proj Projection

	{
		cameras := system.world.GetEntitiesWithComponents(ecs.GetComponentPointerType((*Projection)(nil)))
		if len(cameras) != 1 {
			return projection.ErrWorldShouldHaveOneProjection
		}
		camera := cameras[0]
		if err := system.world.GetComponent(camera, &proj); err != nil {
			return err
		}
		if err := system.world.GetComponent(camera, &cameraTransform); err != nil {
			return err
		}
	}

	var ray shapes.Ray

	{
		mousePos := system.window.NormalizeMouseClick(int(args.X), int(args.Y))
		ray = proj.ShootRay(cameraTransform, mousePos)
	}

	{
		entity := system.getEntity()
		var trans transform.Transform
		if err := system.world.GetComponent(entity, &trans); err != nil {
			return err
		}

		pos := ray.Pos
		// moves ray so it isn't centered
		pos = pos.Add(ray.Rotation().Rotate(camera.Forward).Mul(trans.Size.Z() / 2))
		// moves in front of camera
		pos = pos.Add(cameraTransform.Rotation.Rotate(camera.Forward).Mul(system.distFromCamera))

		system.world.SaveComponent(entity, trans.
			SetPos(pos).
			SetRotation(ray.Rotation()),
		)
	}

	return nil
}
