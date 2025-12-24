package colliderpkg

import (
	"engine/modules/collider"
	"engine/modules/collider/internal/collisions"
	"engine/services/assets"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg { return pkg{} }

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) collider.ToolFactory {
		return collisions.NewToolFactory(
			ioc.Get[logger.Logger](c),
			ioc.Get[assets.Assets](c),
			100,
		)
	})
}
