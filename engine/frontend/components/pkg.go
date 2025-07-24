package components

import "github.com/ogiusek/ioc/v2"

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	// ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, s ecs.)  {})
}
