package mouse

import (
	"frontend/modules/camera"
	"frontend/modules/collider"
	"frontend/modules/transform"
	"frontend/services/assets"
	"frontend/services/media/window"
	"shared/services/ecs"

	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type ShootRayEvent struct{}

func NewShootRayEvent() ShootRayEvent {
	return ShootRayEvent{}
}

type RayChangedTargetEvent struct {
	Camera   ecs.EntityID
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
		events.Listen(w.EventsBuilder(), func(sdl.MouseButtonEvent) {
			events.Emit(s.world.Events(), ShootRayEvent{})
		})

		return nil
	})
}

func (s *cameraRaySystem) Listen(args ShootRayEvent) error {
	mousePos := s.window.NormalizeMousePos(s.window.GetMousePos())

	var nearestCollision collider.ObjectRayCollision
	var nearestCamera ecs.EntityID
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
			nearestCamera = cameraEntity
			continue
		}

		if nearestCollision.Hit().Distance > collision.Hit().Distance {
			nearestCollision = collision
			nearestCamera = cameraEntity
		}
	}

	if nearestCollision == nil {
		if s.hoversOverEntity != nil {
			s.hoversOverEntity = nil
			event := RayChangedTargetEvent{Camera: nearestCamera, EntityID: nil}
			events.Emit(s.events, event)
		}
		return nil
	}

	entity := nearestCollision.Entity()
	if s.hoversOverEntity != nil && *s.hoversOverEntity == entity {
		return nil
	}

	s.hoversOverEntity = &entity
	event := RayChangedTargetEvent{Camera: nearestCamera, EntityID: &entity}
	events.Emit(s.events, event)

	return nil
}
