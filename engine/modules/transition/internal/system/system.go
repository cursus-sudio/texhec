package system

import (
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/logger"
	"slices"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type system struct {
	World         ecs.World      `inject:"1"`
	Logger        logger.Logger  `inject:"1"`
	EventsBuilder events.Builder `inject:"1"`
	Events        events.Events  `inject:"1"`

	delayed []*transition.DelayedEvent
}

func NewSystem(c ioc.Dic) transition.System {
	s := ioc.GetServices[*system](c)

	return ecs.NewSystemRegister(func() error {
		events.Listen(s.EventsBuilder, s.ListenDelayed)
		events.Listen(s.EventsBuilder, s.ListenFrame)
		return nil
	})
}

func (s *system) ListenDelayed(e transition.DelayedEvent) {
	insIdx, _ := slices.BinarySearchFunc(s.delayed, &e, func(a, b *transition.DelayedEvent) int {
		return int(a.Duration - b.Duration)
	})

	s.delayed = slices.Insert(s.delayed, insIdx, &e)
}

func (s *system) ListenFrame(e frames.FrameEvent) {
	toOld := 0
	for _, event := range s.delayed {
		event.Duration -= e.Delta
		if event.Duration > 0 {
			continue
		}

		events.EmitAny(s.Events, event.Event)
		toOld++
	}

	s.delayed = s.delayed[toOld:]
}
