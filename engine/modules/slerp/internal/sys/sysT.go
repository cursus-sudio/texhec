package sys

import (
	"engine/modules/slerp"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
)

type SlerpFn[Component any] func(c1, c2 Component, progress slerp.Progress) Component

type sysT[Component any] struct {
	world          slerp.World
	slerpArray     ecs.ComponentsArray[slerp.SlerpComponent[Component]]
	componentArray ecs.ComponentsArray[Component]
	dirtySet       ecs.DirtySet

	slerpFn SlerpFn[Component]
}

func NewSysT[Component any](
	slerpFn SlerpFn[Component],
) slerp.System {
	return ecs.NewSystemRegister(func(world slerp.World) error {
		s := sysT[Component]{
			world,
			ecs.GetComponentsArray[slerp.SlerpComponent[Component]](world),
			ecs.GetComponentsArray[Component](world),
			ecs.NewDirtySet(),
			slerpFn,
		}

		s.slerpArray.AddDirtySet(s.dirtySet)
		events.Listen(world.EventsBuilder(), s.Listen)

		return nil
	})
}

func (s sysT[Component]) Listen(event frames.FrameEvent) {
	ei := s.dirtySet.Get()

	for _, entity := range ei {
		slerpComponent, ok := s.slerpArray.Get(entity)
		if !ok {
			continue
		}

		slerpComponent.Progress = min(
			slerpComponent.Duration,
			slerpComponent.Progress+event.Delta,
		)
		progress := slerp.Progress(slerpComponent.Progress) / slerp.Progress(slerpComponent.Duration)
		slerpedComponent := s.slerpFn(slerpComponent.From, slerpComponent.To, progress)

		s.slerpArray.Set(entity, slerpComponent)
		s.componentArray.Set(entity, slerpedComponent)
	}
}
