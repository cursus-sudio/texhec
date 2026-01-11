package internal

import (
	"engine/modules/smooth"
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

func NewLastSystem[Component transition.Lerp[Component]](c ioc.Dic) smooth.StopSystem {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*system[Component]](c)
		events.Listen(s.EventsBuilder, func(tick frames.TickEvent) {
			r, ok := s.Record.Entity().Stop(s.Service.recordingID)
			if !ok {
				return
			}
			for _, entity := range r.Entities.GetIndices() {
				beforeComponents, ok := r.Entities.Get(entity)
				if !ok || beforeComponents == nil {
					continue
				}
				before, ok := beforeComponents[0].(Component)
				if !ok {
					continue
				}
				after, ok := s.Service.componentArray.Get(entity)
				if !ok {
					continue
				}
				lerpComponent := transition.NewTransition(before, after, tick.Delta)
				s.Service.lerpArray.Set(entity, lerpComponent)
			}
		})

		return nil
	})
}
