package window

import (
	"frontend/services/frames"
	"shared/services/runtime"

	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type Pkg struct {
	window   *sdl.Window
	renderer *sdl.Renderer
}

func Package(
	window *sdl.Window,
	renderer *sdl.Renderer,
) Pkg {
	return Pkg{
		window:   window,
		renderer: renderer,
	}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Api {
		return newApi(
			pkg.window,
			pkg.renderer,
		)
	})

	ioc.WrapService(b, frames.Draw, func(c ioc.Dic, b frames.Builder) frames.Builder {
		w := ioc.Get[Api](c)
		b.OnFrame(func(of frames.OnFrame) {
			w.Renderer().Present()
			w.Renderer().Clear()
		})
		return b
	})

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b runtime.Builder) runtime.Builder {
		b.OnStop(func() {
			api := ioc.Get[Api](c)
			api.Window().Destroy()
			api.Renderer().Destroy()
		})
		return b
	})
}
