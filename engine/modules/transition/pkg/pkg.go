package transitionpkg

import (
	"engine/modules/transition"
	"engine/modules/transition/internal/easing"
	"engine/modules/transition/internal/tool"
	"engine/modules/transition/internal/transitionimpl"
	"engine/services/codec"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// components
			Register(transition.EasingComponent{}).
			Register(transition.Progress(0))
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) transitionimpl.Builder {
		return transitionimpl.NewBuilder()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) transition.System {
		return ioc.Get[transitionimpl.Builder](c).Build()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) transition.ToolFactory {
		return tool.NewToolFactory()
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) transition.EasingService {
		return easing.NewEasingService()
	})
}
