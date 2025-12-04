package animationpkg

import (
	"engine/modules/animation"
	"engine/modules/animation/internal"
	"engine/services/codec"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			// types
			Register(animation.AnimationState(0)).
			Register(animation.EasingFunctionID(0)).
			Register(animation.Transition{}).
			Register(animation.Event{}).
			Register(animation.AnimationID(0)).
			Register(animation.Animation{}).
			// components
			Register(animation.AnimationComponent{})
	})

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
