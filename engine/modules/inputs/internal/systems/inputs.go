package systems

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/frames"
	mediainputs "engine/services/media/inputs"

	"github.com/ogiusek/events"
)

type inputsSystem struct {
	inputs mediainputs.Api
}

func NewInputsSystem(
	mediainputs mediainputs.Api,
) ecs.SystemRegister[inputs.World] {
	return ecs.NewSystemRegister(func(w inputs.World) error {
		s := &inputsSystem{inputs: mediainputs}
		events.Listen(w.EventsBuilder(), s.Listen)
		return nil
	})
}

func (system *inputsSystem) Listen(args frames.FrameEvent) {
	system.inputs.Poll()
}
