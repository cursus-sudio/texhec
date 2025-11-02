package systems

import (
	"frontend/services/frames"
	"frontend/services/media/inputs"
	"shared/services/ecs"

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
