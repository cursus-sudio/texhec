package colliderpkg

import (
	"engine/modules/collider"
	"engine/modules/collider/internal/collisions"
	"engine/modules/groups"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg { return pkg{} }

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) collider.Service {
		return collisions.NewService(
			ioc.Get[ecs.World](c),
			ioc.Get[transform.Service](c),
			ioc.Get[groups.Service](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[assets.Assets](c),
			100,
		)
	})
}
