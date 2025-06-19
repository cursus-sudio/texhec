package ecs

import "github.com/ogiusek/ioc/v2"

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) WorldFactory { return func() World { return newWorld() } })
}
