package broadcollision

import (
	"frontend/services/assets"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg { return Pkg{} }

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[CollisionService] {
		return factory(
			ioc.Get[assets.Assets](c),
			ioc.Get[logger.Logger](c),
		)
	})
}
