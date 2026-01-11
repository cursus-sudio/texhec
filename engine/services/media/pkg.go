package media

import (
	"engine/services/media/audio"
	"engine/services/media/inputs"
	"engine/services/media/window"
	"engine/services/runtime"

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

	ioc.WrapServiceInOrder(b, runtime.OrderCleanUp, func(c ioc.Dic, b runtime.Builder) {
		b.OnStop(func(r runtime.Runtime) {
			sdl.Quit()
		})
	})
}
