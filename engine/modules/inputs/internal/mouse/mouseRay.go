package mouse

import (
	"engine/modules/inputs"
	"engine/modules/inputs/internal/tool"
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
	inputs.World
	logger logger.Logger
	window window.Api

	targets []inputs.Target
}

func NewCameraRaySystem(
	logger logger.Logger,
	window window.Api,
) inputs.System {
	return ecs.NewSystemRegister(func(w inputs.World) error {
		s := &cameraRaySystem{
			World:  w,
			logger: logger,
			window: window,

			targets: nil,
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

	targets := []inputs.Target{}
	for _, cameraEntity := range s.Camera().Component().GetEntities() {
		ray := s.Camera().ShootRay(cameraEntity, mousePos)

		cameraCollisions := s.Collider().RaycastAll(ray)
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
	events.Emit(s.Events(), tool.RayChangedTargetEvent{
		Targets: targetsCopy,
	})

	return nil
}
