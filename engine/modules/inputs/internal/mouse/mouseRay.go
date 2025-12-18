package mouse

import (
	"engine/modules/collider"
	"engine/modules/inputs"
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
	inputs.World
	logger logger.Logger
	window window.Api
	events events.Events
	assets assets.Assets

	hoversOverEntity *ecs.EntityID
}

func NewCameraRaySystem(
	logger logger.Logger,
	window window.Api,
) ecs.SystemRegister[inputs.World] {
	return ecs.NewSystemRegister(func(w inputs.World) error {
		s := &cameraRaySystem{
			World:  w,
			logger: logger,
			window: window,
			events: w.Events(),

			hoversOverEntity: nil,
		}
		events.ListenE(w.EventsBuilder(), s.Listen)
		events.Listen(w.EventsBuilder(), func(sdl.MouseButtonEvent) {
			events.Emit(s.Events(), ShootRayEvent{})
		})

		return nil
	})
}

func (s *cameraRaySystem) Listen(args ShootRayEvent) error {
	mousePos := s.window.GetMousePos()

	var nearestCollision collider.ObjectRayCollision
	var nearestCamera ecs.EntityID
	for _, cameraEntity := range s.Camera().Component().GetEntities() {
		camera, err := s.Camera().GetObject(cameraEntity)
		if err != nil {
			return err
		}

		ray := camera.ShootRay(mousePos)

		collision, err := s.Collider().ShootRay(ray)
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
