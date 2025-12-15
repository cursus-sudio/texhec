package connectionpkg

import (
	"engine/modules/connection"
	"engine/modules/connection/internal"
	"engine/services/codec"
	"engine/services/ecs"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) connection.System {
		return ecs.NewSystemRegister(func(w ecs.World) error {
			ioc.Get[ecs.ToolFactory[connection.Connection]](c).
				Build(w)
			return nil
		})
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[connection.Connection] {
		return internal.NewToolFactory(
			ioc.Get[codec.Codec](c),
			ioc.Get[logger.Logger](c),
		)
	})
}
