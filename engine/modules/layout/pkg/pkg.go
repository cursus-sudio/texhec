package layoutpkg

import (
	"engine/modules/hierarchy"
	"engine/modules/layout"
	"engine/modules/layout/internal/service"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) layout.Service {
		return service.NewLayoutService(
			ioc.Get[logger.Logger](c),
			ioc.Get[ecs.World](c),
			ioc.Get[hierarchy.Service](c),
			ioc.Get[transform.Service](c),
		)
	})
}
