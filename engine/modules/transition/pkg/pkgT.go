package transitionpkg

import (
	"engine/modules/transition"
	"engine/modules/transition/internal/transitionimpl"
	"engine/services/codec"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkgT[Component transition.Lerp[Component]] struct {
}

func PackageT[Component transition.Lerp[Component]]() ioc.Pkg {
	return pkgT[Component]{}
}

func (pkgT[Component]) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.
			// components
			Register(transition.TransitionComponent[Component]{}).
			// events
			Register(transition.TransitionEvent[Component]{})
	})
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b transitionimpl.Builder) transitionimpl.Builder {
		sys := transitionimpl.NewSysT[Component](
			ioc.Get[logger.Logger](c),
			ioc.Get[transition.EasingService](c),
		)
		b.Register(sys)
		return b
	})
}
