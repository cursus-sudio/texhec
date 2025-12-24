package transitionimpl

import (
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"

	"github.com/ogiusek/events"
)

type sysT[Component transition.Lerp[Component]] struct {
	world    transition.World
	dirtySet ecs.DirtySet

	transitionArray ecs.ComponentsArray[transition.TransitionComponent[Component]]
	easingArray     ecs.ComponentsArray[transition.EasingComponent]
	componentArray  ecs.ComponentsArray[Component]

	logger        logger.Logger
	easingService transition.EasingService
}

func NewSysT[Component transition.Lerp[Component]](
	logger logger.Logger,
	easingService transition.EasingService,
) transition.System {
	return ecs.NewSystemRegister(func(world transition.World) error {
		s := sysT[Component]{
			world,
			ecs.NewDirtySet(),

			ecs.GetComponentsArray[transition.TransitionComponent[Component]](world),
			ecs.GetComponentsArray[transition.EasingComponent](world),
			ecs.GetComponentsArray[Component](world),

			logger,
			easingService,
		}

		events.Listen(world.EventsBuilder(), s.ListenTransition)

		s.transitionArray.AddDirtySet(s.dirtySet)
		events.Listen(world.EventsBuilder(), s.ListenFrame)

		return nil
	})
}

func (s sysT[Component]) ListenTransition(event transition.TransitionEvent[Component]) {
	s.transitionArray.Set(event.Entity, event.Component)
}

func (s sysT[Component]) ListenFrame(event frames.FrameEvent) {
	ei := s.dirtySet.Get()

	for _, entity := range ei {
		transitionComponent, ok := s.transitionArray.Get(entity)
		if !ok {
			continue
		}

		transitionComponent.Progress = min(
			transitionComponent.Duration,
			transitionComponent.Progress+event.Delta,
		)
		progress := transition.Progress(transitionComponent.Progress) / transition.Progress(transitionComponent.Duration)

		easingComponent, ok := s.easingArray.Get(entity)
		if ok {
			if fn, ok := s.easingService.Get(easingComponent.ID); ok {
				progress = fn(progress)
			}
		}

		component := transitionComponent.From.Lerp(transitionComponent.To, float32(progress))

		s.transitionArray.Set(entity, transitionComponent)
		s.componentArray.Set(entity, component)
	}
}
