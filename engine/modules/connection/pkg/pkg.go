package connectionpkg

import (
	"engine/modules/connection"
	"engine/modules/connection/internal"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) connection.System {
		return ecs.NewSystemRegister(func() error {
			ioc.Get[connection.Service](c)
			return nil
		})
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) connection.Service {
		return internal.NewService(
			ioc.Get[codec.Codec](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[ecs.World](c),
			ioc.Get[events.Builder](c),
		)
	})
}
