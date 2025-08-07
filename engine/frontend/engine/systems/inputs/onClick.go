package inputs

import (
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/services/colliders/shapes"
	"frontend/services/ecs"
	"frontend/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/veandco/go-sdl2/sdl"
)

type OnClickSystem[Projection projection.Projection] struct {
	world          ecs.World
	window         window.Api
	distFromCamera float32
	getEntity      func() ecs.EntityId
}

func NewOnClickSystem[Projection projection.Projection](
	world ecs.World,
	window window.Api,
	distFromCamera float32,
	getEntity func() ecs.EntityId,
) OnClickSystem[Projection] {
	return OnClickSystem[Projection]{
		world:          world,
		window:         window,
		distFromCamera: distFromCamera,
		getEntity:      func() ecs.EntityId { return getEntity() },
	}
}

func (system *OnClickSystem[Projection]) Listen(args sdl.MouseMotionEvent) error {
	var ray shapes.Ray
	var cameraTransform transform.Transform

	{
		cameras := system.world.GetEntitiesWithComponents(ecs.GetComponentPointerType((*Projection)(nil)))
		if len(cameras) != 1 {
			return projection.ErrWorldShouldHaveOneProjection
		}
		camera := cameras[0]
		var proj Projection
		if err := system.world.GetComponent(camera, &proj); err != nil {
			return err
		}
		if err := system.world.GetComponent(camera, &cameraTransform); err != nil {
			return err
		}
		mousePos := mgl32.Vec2{float32(args.X), float32(args.Y)}
		w, h := system.window.Window().GetSize()

		mousePos = mgl32.Vec2{
			(2*float32(args.X)/float32(w) - 1),
			-(2*float32(args.Y)/float32(h) - 1),
		}
		ray = proj.ShootRay(cameraTransform, mousePos)
	}

	{
		entity := system.getEntity()
		var trans transform.Transform
		if err := system.world.GetComponent(entity, &trans); err != nil {
			return err
		}
		rotation := ray.Rotation
		trans = trans.
			SetPos(ray.Pos.
				Add(rotation.Rotate(transform.Up).Mul(trans.Size.Y() / 2)).
				// Add(rotation.Rotate(transform.Fo).Mul(trans.Size.Y() / 2)).
				// Add(rotation.Rotate(transform.Up).Mul(trans.Size.Z() / 2)).
				Add(cameraTransform.Rotation.Rotate(transform.Foward).Mul(system.distFromCamera)).
				// Add(cameraTransform.Rotation.Rotate(transform.Up).Mul(system.distFromCamera)).
				Add(mgl32.Vec3{0, 0, 0}),
			).
			SetRotation(rotation)
		system.world.SaveComponent(entity, trans)
	}

	return nil
}
