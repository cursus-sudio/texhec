package eventspkg

import (
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) events.Builder {
		b := events.NewBuilder()
		b.GoroutinePerListener(false)
		return b
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) events.Events { return ioc.Get[events.Builder](c).Build() })
	ioc.RegisterDependency[events.Events, events.Builder](b)
}
