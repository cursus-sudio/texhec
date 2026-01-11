package internal

import (
	"engine/modules/record"
	"engine/modules/smooth"
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
)

func NewLastSystem[Component transition.Lerp[Component]](
	eventsBuilder events.Builder,
	world ecs.World,
	record record.Service,
	service *Service[Component],
) smooth.StopSystem {
	return ecs.NewSystemRegister(func() error {
		events.Listen(eventsBuilder, func(tick frames.TickEvent) {
			r, ok := record.Entity().Stop(service.recordingID)
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
				after, ok := service.componentArray.Get(entity)
				if !ok {
					continue
				}
				lerpComponent := transition.NewTransition(before, after, tick.Delta)
				service.lerpArray.Set(entity, lerpComponent)
			}
		})

		return nil
	})
}
