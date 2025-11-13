package texture

import "github.com/ogiusek/ioc/v2"

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Factory {
		return &factory{}
	})
}
