package media

import (
	"frontend/services/media/audio"
	"frontend/services/media/inputs"
	"frontend/services/media/window"
	"shared/services/runtime"

	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	w *sdl.Window,
	ctx sdl.GLContext,
) pkg {
	return pkg{
		pkgs: []ioc.Pkg{
			audio.Package(),
			inputs.Package(),
			window.Package(w, ctx),
		},
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}

	ioc.RegisterSingleton(b, func(c ioc.Dic) Api {
		return newApi(
			ioc.Get[inputs.Api](c),
			ioc.Get[window.Api](c),
			ioc.Get[audio.Api](c),
		)
	})
	ioc.RegisterDependency[Api, inputs.Api](b)
	ioc.RegisterDependency[Api, window.Api](b)
	ioc.RegisterDependency[Api, audio.Api](b)

	ioc.WrapService(b, runtime.OrderCleanUp, func(c ioc.Dic, b runtime.Builder) runtime.Builder {
		b.OnStop(func(r runtime.Runtime) {
			sdl.Quit()
		})
		return b
	})
}
