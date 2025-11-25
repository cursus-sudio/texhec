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
			ioc.Get[ecs.ToolFactory[transform.Tool]](c),
			ioc.Get[ecs.ToolFactory[collisions.CollisionService]](c),
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[collisions.CollisionService] {
		return collisions.Factory(
			ioc.Get[assets.Assets](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[ecs.ToolFactory[transform.Tool]](c),
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[collider.CollisionTool] {
		return ecs.NewToolFactory(func(w ecs.World) collider.CollisionTool {
			s := ioc.Get[ecs.ToolFactory[collisions.CollisionService]](c)
			return s.Build(w)
		})
	})
}
