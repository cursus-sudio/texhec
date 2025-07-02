package media

import (
	"frontend/services/media/audio"
	"frontend/services/media/inputs"
	"frontend/services/media/window"
	"shared/services/runtime"

	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type Pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	windowPkg window.Pkg,
) Pkg {
	return Pkg{
		pkgs: []ioc.Pkg{
			audio.Package(),
			inputs.Package(),
			windowPkg,
		},
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
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

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b runtime.Builder) runtime.Builder {
		b.OnStop(func() {
			sdl.Quit()
		})
		return b
	})
}
