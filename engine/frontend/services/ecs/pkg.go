package ecs

import (
	"frontend/services/scopes"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterScoped(b, scopes.Scene, func(c ioc.Dic) World { return NewWorld() })
}
