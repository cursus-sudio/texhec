package transformpkg

import (
	"frontend/modules/transform"
	"frontend/modules/transform/internal"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) transform.System {
		return internal.NewPivotPointSystem(ioc.Get[logger.Logger](c))
	})
}
