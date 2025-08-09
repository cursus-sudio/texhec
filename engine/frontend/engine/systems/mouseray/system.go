package mouseray

import (
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/services/colliders"
	"frontend/services/colliders/shapes"
	"frontend/services/ecs"
	"frontend/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/events"
)

type ShootRayEvent[Projection projection.Projection] struct{}

func NewShootRayEvent[Projection projection.Projection]() ShootRayEvent[Projection] {
	return ShootRayEvent[Projection]{}
}

type RayChangedTargetEvent[Projection projection.Projection] struct {
	EntityID *ecs.EntityId
}

type CameraRaySystem[Projection projection.Projection] struct {
	world                  ecs.World
	collider               colliders.ColliderService
	window                 window.Api
	events                 events.Events
	requiredComponentTypes []ecs.ComponentType

	hoversOverEntity *ecs.EntityId
}

func NewCameraRaySystem[Projection projection.Projection](
	world ecs.World,
	collider colliders.ColliderService,
	window window.Api,
	events events.Events,
	requiredComponentTypes []ecs.ComponentType,
) CameraRaySystem[Projection] {
	return CameraRaySystem[Projection]{
		world:                  world,
		collider:               collider,
		window:                 window,
		events:                 events,
		requiredComponentTypes: requiredComponentTypes,
	}
}

func pow2(x float32) float32 { return x * x }
func getDist(x1, x2 mgl32.Vec3) float32 {
	return pow2(x1[0]-x2[0]) + pow2(x1[1]-x2[1]) + pow2(x1[2]-x2[2])
}

type object[Projection projection.Projection] struct {
	Dist     float32
	EntityID ecs.EntityId
}

func (s *CameraRaySystem[Projection]) Listen(args ShootRayEvent[Projection]) error {
	var cameraTransform transform.Transform
	var ray shapes.Ray
	{
		var proj Projection
		cameras := s.world.GetEntitiesWithComponents(
			ecs.GetComponentPointerType((*Projection)(nil)),
		)
		if len(cameras) != 1 {
			return projection.ErrWorldShouldHaveOneProjection
		}
		camera := cameras[0]
		if err := s.world.GetComponents(camera, &proj, &cameraTransform); err != nil {
			return err
		}

		// s.window.
		mousePos := s.window.NormalizeMouseClick(s.window.GetMousePos())
		ray = proj.ShootRay(cameraTransform, mousePos)
	}

	rayCollider := colliders.NewCollider([]colliders.Shape{ray})

	entities := s.world.GetEntitiesWithComponents(
		append(
			s.requiredComponentTypes,
			ecs.GetComponentType(transform.Transform{}),
			ecs.GetComponentType(colliders.Collider{}),
		)...,
	)
	nearestObject := (*object[Projection])(nil)
	for _, entity := range entities {
		var (
			entityTransform transform.Transform
			entityCollider  colliders.Collider
		)
		if err := s.world.GetComponents(entity,
			&entityCollider,
			&entityTransform,
		); err != nil {
			continue
		}
		entityCollider = entityCollider.Apply(entityTransform)

		collision, err := s.collider.Collides(rayCollider, entityCollider)
		if err != nil {
			return err
		}
		if collision == nil {
			continue
		}

		var minDist float32 = -1

		for _, intersection := range collision.Intersections() {
			dist := getDist(cameraTransform.Pos, intersection.PointOnB())
			if minDist < 0 {
				minDist = dist
				continue
			}
			minDist = min(minDist, dist)
		}

		if minDist < 0 {
			minDist = getDist(cameraTransform.Pos, entityTransform.Pos)
		}

		current := object[Projection]{
			Dist:     minDist,
			EntityID: entity,
		}

		if nearestObject == nil {
			nearestObject = &current
			continue
		}

		if current.Dist > nearestObject.Dist {
			continue
		}

		nearestObject = &current

	}

	if nearestObject != nil {
		if s.hoversOverEntity != nil && *s.hoversOverEntity == nearestObject.EntityID {
			return nil
		}

		s.hoversOverEntity = &nearestObject.EntityID
		event := RayChangedTargetEvent[Projection]{EntityID: &nearestObject.EntityID}
		events.Emit(s.events, event)
	} else if s.hoversOverEntity != nil {
		s.hoversOverEntity = nil
		event := RayChangedTargetEvent[Projection]{EntityID: nil}
		events.Emit(s.events, event)
	}

	return nil
}
