package ecs

import (
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

type PersistentWorld World

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) PersistentWorld { return NewWorld() })
}
