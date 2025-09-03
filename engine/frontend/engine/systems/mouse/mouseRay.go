package mouse

import (
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/engine/tools/broadcollision"
	"frontend/engine/tools/worldprojections"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/media/window"

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
	world           ecs.World
	broadCollisions broadcollision.CollisionDetectionService
	window          window.Api
	events          events.Events
	assets          assets.Assets

	hoversOverEntites map[ecs.ComponentType]ecs.EntityID
}

func NewCameraRaySystem(
	world ecs.World,
	collider broadcollision.CollisionDetectionService,
	window window.Api,
	events events.Events,
) CameraRaySystem {
	return CameraRaySystem{
		world:           world,
		broadCollisions: collider,
		window:          window,
		events:          events,

		hoversOverEntites: map[ecs.ComponentType]ecs.EntityID{},
	}
}

func (s *CameraRaySystem) Listen(args ShootRayEvent) error {
	mousePos := s.window.NormalizeMouseClick(s.window.GetMousePos())

	p, err := ecs.GetRegister[worldprojections.WorldProjectionsRegister](s.world)
	if err != nil {
		return err
	}

	for _, projectionType := range p.Projections.Get() {
		if projectionType == ecs.GetComponentType(projection.Perspective{}) {
			continue
		}
		var cameraTransform transform.Transform
		query := s.world.QueryEntitiesWithComponents(projectionType)
		cameras := query.Entities()
		var nearestCollision broadcollision.ObjectRayCollision
		for _, camera := range cameras {
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

			ray := proj.ShootRay(cameraTransform, mousePos)

			collision, err := s.broadCollisions.ShootRay(ray)
			if err != nil {
				return err
			}
			if collision == nil {
				continue
			}
			if nearestCollision == nil {
				nearestCollision = collision
				continue
			}
			if nearestCollision.Hit().Distance > nearestCollision.Hit().Distance {
				nearestCollision = collision
			}
		}
		if nearestCollision == nil {
			if _, ok := s.hoversOverEntites[projectionType]; ok {
				delete(s.hoversOverEntites, projectionType)
				event := RayChangedTargetEvent{ProjectionType: projectionType, EntityID: nil}
				events.Emit(s.events, event)
			}
			continue
		}

		entity := nearestCollision.Entity()
		hoversOverEntity, _ := s.hoversOverEntites[projectionType]
		if hoversOverEntity == entity {
			continue
		}

		s.hoversOverEntites[projectionType] = entity
		event := RayChangedTargetEvent{ProjectionType: projectionType, EntityID: &entity}
		events.Emit(s.events, event)
	}

	return nil
}
