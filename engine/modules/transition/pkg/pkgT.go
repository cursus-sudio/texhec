package transitionpkg

import (
	"engine/modules/transition"
	"engine/modules/transition/internal/transitionimpl"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkgT[Component transition.Lerp[Component]] struct {
}

func PackageT[Component transition.Lerp[Component]]() ioc.Pkg {
	return pkgT[Component]{}
}

func (pkgT[Component]) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// components
			Register(transition.TransitionComponent[Component]{}).
			// events
			Register(transition.TransitionEvent[Component]{})
	})
	ioc.WrapService(b, func(c ioc.Dic, b transitionimpl.Builder) {
		sys := transitionimpl.NewSysT[Component](
			ioc.Get[ecs.World](c),
			ioc.Get[events.Builder](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[transition.EasingService](c),
		)
		b.Register(sys)
	})
}
