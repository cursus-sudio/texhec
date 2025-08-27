package colliders

import (
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
}

func Package() Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ColliderServiceBuilder { return NewBuilder() })

	ioc.RegisterSingleton(b, func(c ioc.Dic) ColliderService {
		s, errs := ioc.Get[ColliderServiceBuilder](c).Build()
		if len(errs) != 0 {
			logger := ioc.Get[logger.Logger](c)
			logger.Fatal(errs...)
		}
		return s
	})
}
