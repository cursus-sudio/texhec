package draw

import "github.com/ogiusek/ioc"

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(c ioc.Dic) {
	ioc.RegisterSingleton(c, func(c ioc.Dic) DrawApi {
		return &drawApi{}
	})
}
