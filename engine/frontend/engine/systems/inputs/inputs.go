package inputssys

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
	return &inputsSystem{inputs}
}

func (s *inputsSystem) Register(b events.Builder) {
	events.Listen(b, s.Listen)
}

func (system *inputsSystem) Listen(args frames.FrameEvent) {
	system.inputs.Poll()
}
