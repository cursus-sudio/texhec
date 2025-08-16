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
	getEntity      func() ecs.EntityID
	shootRay       func(shapes.Ray)
	camerasQuery   ecs.LiveQuery
}

func NewShootRaySystem[Projection projection.Projection](
	world ecs.World,
	window window.Api,
	console console.Console,
	distFromCamera float32,
	getEntity func() ecs.EntityID,
	shootRay func(shapes.Ray),
) ShootRaySystem[Projection] {
	camerasLiveQuery := world.QueryEntitiesWithComponents(ecs.GetComponentPointerType((*Projection)(nil)))
	return ShootRaySystem[Projection]{
		world:          world,
		window:         window,
		console:        console,
		distFromCamera: distFromCamera,
		getEntity:      func() ecs.EntityID { return getEntity() },
		shootRay:       shootRay,
		camerasQuery:   camerasLiveQuery,
	}
}

func (system *ShootRaySystem[Projection]) Listen(args ShootRayEvent) error {
	cameras := system.camerasQuery.Entities()
	if len(cameras) != 1 {
		return projection.ErrWorldShouldHaveOneProjection
	}
	cameraEntity := cameras[0]
	cameraTransform, err := ecs.GetComponent[transform.Transform](system.world, cameraEntity)
	if err != nil {
		return err
	}
	proj, err := ecs.GetComponent[Projection](system.world, cameraEntity)
	if err != nil {
		return err
	}

	var ray shapes.Ray

	{
		mousePos := system.window.NormalizeMouseClick(int(args.X), int(args.Y))
		ray = proj.ShootRay(cameraTransform, mousePos)
		system.shootRay(ray)
	}

	{
		entity := system.getEntity()
		trans, err := ecs.GetComponent[transform.Transform](system.world, entity)
		if err != nil {
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
