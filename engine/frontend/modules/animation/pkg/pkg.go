package animationpkg

import (
	"frontend/modules/animation"
	"frontend/modules/animation/internal"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) internal.AnimationSystemBuilder {
		return internal.NewBuilder(ioc.Get[logger.Logger](c))
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) animation.AnimationSystemBuilder {
		return ioc.Get[internal.AnimationSystemBuilder](c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) animation.System {
		return ioc.Get[internal.AnimationSystemBuilder](c).Build()
	})
}
