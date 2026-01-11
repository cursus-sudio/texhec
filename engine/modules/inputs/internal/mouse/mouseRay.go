package mouse

import (
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/inputs"
	"engine/modules/inputs/internal/service"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"
	"slices"

	"github.com/ogiusek/events"
	"github.com/veandco/go-sdl2/sdl"
)

type ShootRayEvent struct{}

func NewShootRayEvent() ShootRayEvent {
	return ShootRayEvent{}
}

type cameraRaySystem struct {
	world    ecs.World
	camera   camera.Service
	collider collider.Service

	events events.Events
	logger logger.Logger
	window window.Api

	targets []inputs.Target
}

func NewCameraRaySystem(
	eventsBuilder events.Builder,
	world ecs.World,
	camera camera.Service,
	collider collider.Service,
	logger logger.Logger,
	window window.Api,
) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := &cameraRaySystem{
			world:    world,
			camera:   camera,
			collider: collider,

			events: eventsBuilder.Events(),
			logger: logger,
			window: window,

			targets: nil,
		}
		events.ListenE(eventsBuilder, s.Listen)
		events.Listen(eventsBuilder, func(sdl.MouseButtonEvent) {
			events.Emit(eventsBuilder.Events(), ShootRayEvent{})
		})

		return nil
	})
}

func (s *cameraRaySystem) Listen(args ShootRayEvent) error {
	mousePos := s.window.GetMousePos()

	targets := []inputs.Target{}
	for _, cameraEntity := range s.camera.Component().GetEntities() {
		ray := s.camera.ShootRay(cameraEntity, mousePos)

		cameraCollisions := s.collider.RaycastAll(ray)
		for _, collision := range cameraCollisions {
			target := inputs.Target{
				ObjectRayCollision: collision,
				Camera:             cameraEntity,
			}
			targets = append(targets, target)
		}
	}

	slices.SortFunc(targets, func(a, b inputs.Target) int {
		if a.Hit.Distance < b.Hit.Distance {
			return -1
		}
		if a.Hit.Distance > b.Hit.Distance {
			return 1
		}
		return 0
	})

	if slices.Equal(s.targets, targets) {
		return nil
	}

	s.targets = targets

	targetsCopy := make([]inputs.Target, len(s.targets))
	copy(targetsCopy, s.targets)
	events.Emit(s.events, service.RayChangedTargetEvent{
		Targets: targetsCopy,
	})

	return nil
}
