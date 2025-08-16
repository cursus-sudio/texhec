package mouse

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

type ShootRayEvent struct{}

func NewShootRayEvent() ShootRayEvent {
	return ShootRayEvent{}
}

type RayChangedTargetEvent struct {
	ProjectionType ecs.ComponentType
	EntityID       *ecs.EntityID
}

type CameraRaySystem struct {
	world    ecs.World
	collider colliders.ColliderService
	window   window.Api
	events   events.Events

	liveQuery ecs.LiveQuery
	// requiredComponentTypes []ecs.ComponentType

	projectionTypes   []ecs.ComponentType
	hoversOverEntites map[ecs.ComponentType]ecs.EntityID
}

func NewCameraRaySystem(
	world ecs.World,
	collider colliders.ColliderService,
	window window.Api,
	events events.Events,
	cameraProjections []ecs.ComponentType,
	requiredComponentTypes []ecs.ComponentType,
) CameraRaySystem {
	liveQuery := world.QueryEntitiesWithComponents(
		append(
			requiredComponentTypes,
			ecs.GetComponentType(transform.Transform{}),
			ecs.GetComponentType(colliders.Collider{}),
			ecs.GetComponentType(projection.UsedProjection{}),
		)...,
	)

	return CameraRaySystem{
		world:    world,
		collider: collider,
		window:   window,
		events:   events,

		liveQuery: liveQuery,

		projectionTypes:   cameraProjections,
		hoversOverEntites: map[ecs.ComponentType]ecs.EntityID{},
	}
}

func pow2(x float32) float32 { return x * x }
func getDist(x1, x2 mgl32.Vec3) float32 {
	return pow2(x1[0]-x2[0]) + pow2(x1[1]-x2[1]) + pow2(x1[2]-x2[2])
}

type object struct {
	Dist     float32
	EntityID ecs.EntityID
	ok       bool
}

func (obj object) Ok() bool {
	return obj.ok
}

func (s *CameraRaySystem) Listen(args ShootRayEvent) error {
	mousePos := s.window.NormalizeMouseClick(s.window.GetMousePos())

	cameraTransforms := make(map[ecs.ComponentType]transform.Transform, len(s.projectionTypes))
	rays := make(map[ecs.ComponentType]colliders.Collider, len(s.projectionTypes))

	for _, projectionType := range s.projectionTypes {
		var cameraTransform transform.Transform
		var ray shapes.Ray
		query := s.world.QueryEntitiesWithComponents(projectionType)
		cameras := query.Entities()
		if len(cameras) != 1 {
			return projection.ErrWorldShouldHaveOneProjection
		}
		camera := cameras[0]
		var err error
		cameraTransform, err = ecs.GetComponent[transform.Transform](s.world, camera)
		if err != nil {
			return err
		}
		anyProj, err := s.world.GetComponent(camera, projectionType)
		if err != nil {
			return err
		}
		proj, ok := anyProj.(projection.Projection)
		if !ok {
			return projection.ErrExpectedUsedProjectionToImplementProjection
		}

		ray = proj.ShootRay(cameraTransform, mousePos)

		rayCollider := colliders.NewCollider([]colliders.Shape{ray})
		rays[projectionType] = rayCollider
		cameraTransforms[projectionType] = cameraTransform
	}

	nearestObjects := make(map[ecs.ComponentType]*object)

	for _, projectionType := range s.projectionTypes {
		nearestObjects[projectionType] = &object{}
	}

	for _, entity := range s.liveQuery.Entities() {
		usedProjection, err := ecs.GetComponent[projection.UsedProjection](s.world, entity)
		if err != nil {
			continue
		}
		rayCollider, ok := rays[usedProjection.ProjectionComponent]
		if !ok {
			continue
		}
		cameraTransform, ok := cameraTransforms[usedProjection.ProjectionComponent]
		if !ok {
			continue
		}
		nearestObject, ok := nearestObjects[usedProjection.ProjectionComponent]
		if !ok {
			continue
		}
		entityTransform, err := ecs.GetComponent[transform.Transform](s.world, entity)
		if err != nil {
			continue
		}
		entityCollider, err := ecs.GetComponent[colliders.Collider](s.world, entity)
		if err != nil {
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

		current := object{
			Dist:     minDist,
			EntityID: entity,
			ok:       true,
		}

		if !nearestObject.Ok() {
			*nearestObject = current
			continue
		}

		if current.Dist > nearestObject.Dist {
			continue
		}

		*nearestObject = current
	}

	for projectionType, nearestObject := range nearestObjects {
		if nearestObject.Ok() {
			hoversOverEntity, _ := s.hoversOverEntites[projectionType]
			if hoversOverEntity == nearestObject.EntityID {
				continue
			}

			id := nearestObject.EntityID
			s.hoversOverEntites[projectionType] = id
			event := RayChangedTargetEvent{ProjectionType: projectionType, EntityID: &id}
			events.Emit(s.events, event)
		} else if _, ok := s.hoversOverEntites[projectionType]; ok {
			delete(s.hoversOverEntites, projectionType)
			event := RayChangedTargetEvent{ProjectionType: projectionType, EntityID: nil}
			events.Emit(s.events, event)
		}
	}

	return nil
}
