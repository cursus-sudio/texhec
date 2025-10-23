package mousesys

import (
	"frontend/engine/components/projection"
	"frontend/engine/components/transform"
	"frontend/engine/tools/broadcollision"
	"frontend/engine/tools/cameras"
	"frontend/engine/tools/worldprojections"
	"frontend/services/assets"
	"frontend/services/media/window"
	"shared/services/ecs"

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

type cameraRaySystem struct {
	world           ecs.World
	transformArray  ecs.ComponentsArray[transform.Transform]
	broadCollisions broadcollision.CollisionDetectionService
	window          window.Api
	events          events.Events
	assets          assets.Assets
	cameraCtors     cameras.CameraConstructors

	hoversOverEntites map[ecs.ComponentType]ecs.EntityID
}

func NewCameraRaySystem(
	world ecs.World,
	collider broadcollision.CollisionDetectionService,
	window window.Api,
	events events.Events,
	cameraCtors cameras.CameraConstructors,
) ecs.SystemRegister {
	return &cameraRaySystem{
		world:           world,
		transformArray:  ecs.GetComponentsArray[transform.Transform](world.Components()),
		broadCollisions: collider,
		window:          window,
		events:          events,
		cameraCtors:     cameraCtors,

		hoversOverEntites: map[ecs.ComponentType]ecs.EntityID{},
	}
}

func (s *cameraRaySystem) Register(b events.Builder) {
	events.ListenE(b, s.Listen)
}

func (s *cameraRaySystem) Listen(args ShootRayEvent) error {
	mousePos := s.window.NormalizeMousePos(s.window.GetMousePos())

	p, err := ecs.GetRegister[worldprojections.WorldProjectionsRegister](s.world)
	if err != nil {
		return err
	}

	for _, projectionType := range p.Projections.Get() {
		if projectionType == ecs.GetComponentType(projection.Perspective{}) {
			continue
		}
		query := s.world.QueryEntitiesWithComponents(projectionType)
		cameras := query.Entities()
		var nearestCollision broadcollision.ObjectRayCollision
		for _, cameraEntity := range cameras {
			camera, err := s.cameraCtors.Get(cameraEntity, projectionType)
			if err != nil {
				return err
			}

			ray := camera.ShootRay(mousePos)

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

			if nearestCollision.Hit().Distance > collision.Hit().Distance {
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
