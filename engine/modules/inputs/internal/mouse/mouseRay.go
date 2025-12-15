package mouse

import (
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/services/assets"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"

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
	logger          logger.Logger
	cameraArray     ecs.ComponentsArray[camera.CameraComponent]
	broadCollisions collider.Interface
	window          window.Api
	events          events.Events
	assets          assets.Assets
	cameraResolver  camera.Interface

	hoversOverEntity *ecs.EntityID
}

func NewCameraRaySystem(
	logger logger.Logger,
	colliderFactory ecs.ToolFactory[collider.Collider],
	window window.Api,
	cameraResolver ecs.ToolFactory[camera.Camera],
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &cameraRaySystem{
			world:           w,
			logger:          logger,
			cameraArray:     ecs.GetComponentsArray[camera.CameraComponent](w),
			broadCollisions: colliderFactory.Build(w).Collider(),
			window:          window,
			events:          w.Events(),
			cameraResolver:  cameraResolver.Build(w).Camera(),

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
	mousePos := s.window.GetMousePos()

	var nearestCollision collider.ObjectRayCollision
	var nearestCamera ecs.EntityID
	for _, cameraEntity := range s.cameraArray.GetEntities() {
		camera, err := s.cameraResolver.GetObject(cameraEntity)
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
