package internal

import (
	"engine/modules/record"
	"engine/modules/transition"
	"engine/services/ecs"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type Service[Component transition.Lerp[Component]] struct {
	World       ecs.World `inject:"1"`
	recordingID record.RecordingID
	config      record.Config

	componentArray ecs.ComponentsArray[Component]
	lerpArray      ecs.ComponentsArray[transition.TransitionComponent[Component]]
}

func NewService[Component transition.Lerp[Component]](c ioc.Dic) *Service[Component] {
	config := record.NewConfig()
	record.AddToConfig[Component](config)

	s := ioc.GetServices[*Service[Component]](c)

	s.recordingID = 0
	s.config = config
	s.componentArray = ecs.GetComponentsArray[Component](s.World)
	s.lerpArray = ecs.GetComponentsArray[transition.TransitionComponent[Component]](s.World)
	return s
}

//

type system[Component transition.Lerp[Component]] struct {
	EventsBuilder events.Builder      `inject:"1"`
	World         ecs.World           `inject:"1"`
	Record        record.Service      `inject:"1"`
	Service       *Service[Component] `inject:"1"`
}
