package renderpkg

import (
	"engine/modules/render"
	"engine/modules/render/internal"
	transitionpkg "engine/modules/transition/pkg"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	for _, pkg := range []ioc.Pkg{
		transitionpkg.PackageT[render.ColorComponent](),
		transitionpkg.PackageT[render.TextureFrameComponent](),
	} {
		pkg.Register(b)
	}

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.Service {
		return internal.NewTool(
			ioc.Get[ecs.World](c),
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.System {
		return ecs.NewSystemRegister(func() error {
			ecs.RegisterSystems(
				internal.NewClearSystem(ioc.Get[events.Builder](c)),
				internal.NewErrorLogger(
					ioc.Get[logger.Logger](c),
					ioc.Get[render.Service](c),
					ioc.Get[events.Builder](c),
				),
				internal.NewRenderSystem(
					ioc.Get[ecs.World](c),
					ioc.Get[window.Api](c),
					ioc.Get[events.Builder](c),
				),
			)
			return nil
		})
	})
}
