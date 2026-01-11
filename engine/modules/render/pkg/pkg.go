package renderpkg

import (
	"engine/modules/render"
	"engine/modules/render/internal"
	transitionpkg "engine/modules/transition/pkg"
	"engine/services/ecs"

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
		return internal.NewService(c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) render.System {
		return ecs.NewSystemRegister(func() error {
			ecs.RegisterSystems(
				internal.NewClearSystem(c),
				internal.NewErrorLogger(c),
				internal.NewRenderSystem(c),
			)
			return nil
		})
	})
}
