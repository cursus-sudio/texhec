package ecs

import (
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

// package registers example singleton world
func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) World { return NewWorld() })
}
