package inputs

import (
	"frontend/services/frames"
	"shared/services/clock"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) *inputsApi {
		return newInputsApi(
			ioc.Get[clock.Clock](c),
			ioc.Get[events.Events](c),
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) InputsApi { return ioc.Get[*inputsApi](c) })

	// ioc.RegisterSingleton()
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b frames.Builder) frames.Builder {
		i := ioc.Get[inputsApi](c)
		return b.OnFrame(func(of frames.OnFrame) {
			i.Poll()
		})
	})
}
