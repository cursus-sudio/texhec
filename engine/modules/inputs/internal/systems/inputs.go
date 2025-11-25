package systems

import (
	"engine/services/ecs"
	"engine/services/frames"
	"engine/services/media/inputs"

	"github.com/ogiusek/events"
)

type inputsSystem struct {
	inputs inputs.Api
}

func NewInputsSystem(
	inputs inputs.Api,
) ecs.SystemRegister {
	return ecs.NewSystemRegister(func(w ecs.World) error {
		s := &inputsSystem{inputs: inputs}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (system *inputsSystem) Listen(args frames.FrameEvent) {
	system.inputs.Poll()
}
