package mouse

import (
	"frontend/engine/camera"
	"frontend/engine/collider"
	"frontend/engine/transform"
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
	EntityID *ecs.EntityID
}

type cameraRaySystem struct {
	world           ecs.World
	transformArray  ecs.ComponentsArray[transform.TransformComponent]
	cameraArray     ecs.ComponentsArray[camera.CameraComponent]
	broadCollisions collider.CollisionTool
	window          window.Api
	events          events.Events
	assets          assets.Assets
	cameraResolver  camera.CameraTool

	hoversOverEntity *ecs.EntityID
}

func NewCameraRaySystem(
	colliderFactory ecs.ToolFactory[collider.CollisionTool],
	window window.Api,
	cameraResolver ecs.ToolFactory[camera.CameraTool],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &cameraRaySystem{
			world:           w,
			transformArray:  ecs.GetComponentsArray[transform.TransformComponent](w.Components()),
			cameraArray:     ecs.GetComponentsArray[camera.CameraComponent](w.Components()),
			broadCollisions: colliderFactory.Build(w),
			window:          window,
			events:          w.Events(),
			cameraResolver:  cameraResolver.Build(w),

			hoversOverEntity: nil,
		}
		events.ListenE(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (s *cameraRaySystem) Listen(args ShootRayEvent) error {
	mousePos := s.window.NormalizeMousePos(s.window.GetMousePos())

	var nearestCollision collider.ObjectRayCollision
	for _, cameraEntity := range s.cameraArray.GetEntities() {
		camera, err := s.cameraResolver.Get(cameraEntity)
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
		if s.hoversOverEntity != nil {
			s.hoversOverEntity = nil
			event := RayChangedTargetEvent{EntityID: nil}
			events.Emit(s.events, event)
		}
		return nil
	}

	entity := nearestCollision.Entity()
	if s.hoversOverEntity != nil && *s.hoversOverEntity == entity {
		return nil
	}

	s.hoversOverEntity = &entity
	event := RayChangedTargetEvent{EntityID: &entity}
	events.Emit(s.events, event)

	return nil
}
