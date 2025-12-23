package renderpkg

import (
	"engine/modules/render"
	"engine/modules/render/internal"
	transitionpkg "engine/modules/transition/pkg"
	"engine/services/ecs"
	"engine/services/logger"
	"engine/services/media/window"

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

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.ToolFactory {
		return internal.NewTool()
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.System {
		return ecs.NewSystemRegister(func(w render.World) error {
			ecs.RegisterSystems(w,
				internal.NewClearSystem(),
				internal.NewErrorLogger(
					ioc.Get[logger.Logger](c),
					ioc.Get[render.ToolFactory](c).Build(w),
				),
				internal.NewRenderSystem(ioc.Get[window.Api](c)),
			)
			return nil
		})
	})
}
