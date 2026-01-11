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
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type ShootRayEvent struct{}

func NewShootRayEvent() ShootRayEvent {
	return ShootRayEvent{}
}

type cameraRaySystem struct {
	World    ecs.World        `inject:"1"`
	Camera   camera.Service   `inject:"1"`
	Collider collider.Service `inject:"1"`

	EventsBuilder events.Builder `inject:"1"`
	Events        events.Events  `inject:"1"`
	Logger        logger.Logger  `inject:"1"`
	Window        window.Api     `inject:"1"`

	targets []inputs.Target
}

func NewCameraRaySystem(c ioc.Dic) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*cameraRaySystem](c)
		s.targets = nil

		events.ListenE(s.EventsBuilder, s.Listen)
		events.Listen(s.EventsBuilder, func(sdl.MouseButtonEvent) {
			events.Emit(s.EventsBuilder.Events(), ShootRayEvent{})
		})

		return nil
	})
}

func (s *cameraRaySystem) Listen(args ShootRayEvent) error {
	mousePos := s.Window.GetMousePos()

	targets := []inputs.Target{}
	for _, cameraEntity := range s.Camera.Component().GetEntities() {
		ray := s.Camera.ShootRay(cameraEntity, mousePos)

		cameraCollisions := s.Collider.RaycastAll(ray)
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
	events.Emit(s.Events, service.RayChangedTargetEvent{
		Targets: targetsCopy,
	})

	return nil
}
