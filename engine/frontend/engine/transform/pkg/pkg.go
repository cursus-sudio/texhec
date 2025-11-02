package transformpkg

import (
	"frontend/engine/transform"
	"frontend/engine/transform/internal"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() ioc.Pkg {
	return Pkg{}
}

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) transform.System {
		return internal.NewPivotPointSystem(ioc.Get[logger.Logger](c))
	})
}
