package ecs

import (
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

// package registers example singleton world
func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) World { return NewWorld() })
}
