package eventspkg

import (
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) events.Builder {
		b := events.NewBuilder()
		b.GoroutinePerListener(false)
		return b
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) events.Events { return ioc.Get[events.Builder](c).Events() })
}
