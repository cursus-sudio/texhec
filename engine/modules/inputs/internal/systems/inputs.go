package systems

import (
	"engine/modules/inputs"
	"engine/services/ecs"
	"engine/services/frames"
	mediainputs "engine/services/media/inputs"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type inputsSystem struct {
	Inputs        mediainputs.Api `inject:"1"`
	EventsBuilder events.Builder  `inject:"1"`
}

func NewInputsSystem(c ioc.Dic) inputs.System {
	return ecs.NewSystemRegister(func() error {
		s := ioc.GetServices[*inputsSystem](c)
		events.Listen(s.EventsBuilder, s.Listen)
		return nil
	})
}

func (system *inputsSystem) Listen(args frames.FrameEvent) {
	system.Inputs.Poll()
}
