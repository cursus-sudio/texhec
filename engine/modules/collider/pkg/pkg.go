package colliderpkg

import (
	"engine/modules/collider"
	"engine/modules/collider/internal/collisions"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg { return pkg{} }

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) collider.System {
		return collisions.NewColliderSystem(
			ioc.Get[logger.Logger](c),
			ioc.Get[ecs.ToolFactory[transform.Transform]](c),
			ioc.Get[ecs.ToolFactory[collider.Collider]](c),
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[collider.Collider] {
		return collisions.NewToolFactory(
			ioc.Get[logger.Logger](c),
			ioc.Get[assets.Assets](c),
			ioc.Get[ecs.ToolFactory[transform.Transform]](c),
			100,
		)
	})
}
