package internal

import (
	"engine/modules/record"
	"engine/modules/transition"
	"engine/services/ecs"
)

type state[Component transition.Lerp[Component]] struct {
	recordingID record.RecordingID
	config      record.Config

	componentArray ecs.ComponentsArray[Component]
	lerpArray      ecs.ComponentsArray[transition.TransitionComponent[Component]]
}

type Service[Component transition.Lerp[Component]] struct {
	*state[Component]
}

func NewService[Component transition.Lerp[Component]](w ecs.World) *Service[Component] {
	config := record.NewConfig()
	record.AddToConfig[Component](config)

	return &Service[Component]{
		state: &state[Component]{
			recordingID: 0,
			config:      config,

			componentArray: ecs.GetComponentsArray[Component](w),
			lerpArray:      ecs.GetComponentsArray[transition.TransitionComponent[Component]](w),
		},
	}
}
