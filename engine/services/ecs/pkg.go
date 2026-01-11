package ecs

import (
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

// package registers example singleton world
func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) events.Builder {
		return events.NewBuilder()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) events.Events {
		return ioc.Get[events.Builder](c).Build()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) World {
		return NewWorld()
	})
}
