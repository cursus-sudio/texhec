package internal

import (
	"engine/modules/record"
	"engine/modules/smooth"
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
)

func NewFirstSystem[Component transition.Lerp[Component]](
	eventsBuilder events.Builder,
	world ecs.World,
	record record.Service,
	service *Service[Component],
) smooth.StartSystem {
	return ecs.NewSystemRegister(func() error {
		events.Listen(eventsBuilder, func(tick frames.TickEvent) {
			for _, entity := range service.lerpArray.GetEntities() {
				transitionComponent, ok := service.lerpArray.Get(entity)
				if !ok {
					continue
				}
				service.lerpArray.Remove(entity)
				service.componentArray.Set(entity, transitionComponent.To)
			}

			service.recordingID = record.Entity().StartBackwardsRecording(service.config)
		})

		return nil
	})
}
