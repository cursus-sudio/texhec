package colliderpkg

import (
	"frontend/modules/collider"
	"frontend/modules/collider/internal/collisions"
	"frontend/services/assets"
	"shared/services/ecs"
	"shared/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg { return Pkg{} }

func (Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) collider.System {
		return collisions.NewColliderSystem(
			ioc.Get[ecs.ToolFactory[collisions.CollisionService]](c),
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[collisions.CollisionService] {
		return collisions.Factory(
			ioc.Get[assets.Assets](c),
			ioc.Get[logger.Logger](c),
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[collider.CollisionTool] {
		return ecs.NewToolFactory(func(w ecs.World) collider.CollisionTool {
			s := ioc.Get[ecs.ToolFactory[collisions.CollisionService]](c)
			return s.Build(w)
		})
	})
}
