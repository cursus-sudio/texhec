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
	eventsBuilder events.Builder,
	mediainputs mediainputs.Api,
) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := &inputsSystem{inputs: mediainputs}
		events.Listen(eventsBuilder, s.Listen)
		return nil
	})
}

func (system *inputsSystem) Listen(args frames.FrameEvent) {
	system.inputs.Poll()
}
