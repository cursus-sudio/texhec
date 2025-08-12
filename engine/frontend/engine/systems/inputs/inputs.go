package inputs

import (
	"frontend/services/frames"
	"frontend/services/media/inputs"
)

type InputsSystem struct {
	inputs inputs.Api
}

func NewInputsSystem(
	inputs inputs.Api,
) InputsSystem {
	return InputsSystem{inputs}
}

func (system *InputsSystem) Update(args frames.FrameEvent) {
	system.inputs.Poll()
}
