package internal

import (
	"engine/modules/smooth"
	"engine/modules/transition"
	"engine/services/ecs"
	"engine/services/frames"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

func NewFirstSystem[Component transition.Lerp[Component]](c ioc.Dic) smooth.StartSystem {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*system[Component]](c)
		events.Listen(s.EventsBuilder, func(tick frames.TickEvent) {
			for _, entity := range s.Service.lerpArray.GetEntities() {
				transitionComponent, ok := s.Service.lerpArray.Get(entity)
				if !ok {
					continue
				}
				s.Service.lerpArray.Remove(entity)
				s.Service.componentArray.Set(entity, transitionComponent.To)
			}

			s.Service.recordingID = s.Record.Entity().StartBackwardsRecording(s.Service.config)
		})

		return nil
	})
}
