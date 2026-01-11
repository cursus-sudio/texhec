package transitionimpl

import (
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type sysT[Component transition.Lerp[Component]] struct {
	World         ecs.World                `inject:"1"`
	Logger        logger.Logger            `inject:"1"`
	EasingService transition.EasingService `inject:"1"`
	EventsBuilder events.Builder           `inject:"1"`

	dirtySet ecs.DirtySet

	transitionArray ecs.ComponentsArray[transition.TransitionComponent[Component]]
	easingArray     ecs.ComponentsArray[transition.EasingComponent]
	componentArray  ecs.ComponentsArray[Component]
}

func NewSysT[Component transition.Lerp[Component]](c ioc.Dic) transition.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*sysT[Component]](c)

		s.dirtySet = ecs.NewDirtySet()
		s.transitionArray = ecs.GetComponentsArray[transition.TransitionComponent[Component]](s.World)
		s.easingArray = ecs.GetComponentsArray[transition.EasingComponent](s.World)
		s.componentArray = ecs.GetComponentsArray[Component](s.World)

		events.Listen(s.EventsBuilder, s.ListenTransition)

		s.transitionArray.AddDirtySet(s.dirtySet)
		events.Listen(s.EventsBuilder, s.ListenFrame)

		return nil
	})
}

func (s *sysT[Component]) ListenTransition(event transition.TransitionEvent[Component]) {
	s.transitionArray.Set(event.Entity, event.Component)
}

func (s *sysT[Component]) ListenFrame(event frames.FrameEvent) {
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
			if fn, ok := s.EasingService.Get(easingComponent.ID); ok {
				progress = fn(progress)
			}
		}

		component := transitionComponent.From.Lerp(transitionComponent.To, float32(progress))

		s.transitionArray.Set(entity, transitionComponent)
		s.componentArray.Set(entity, component)
	}
}
