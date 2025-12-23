package window

import (
	"engine/services/logger"
	runtimeservice "engine/services/runtime"

	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type pkg struct {
	window  *sdl.Window
	context sdl.GLContext
}

func Package(
	window *sdl.Window,
	context sdl.GLContext,
) ioc.Pkg {
	return pkg{
		window:  window,
		context: context,
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Api {
		return newApi(
			pkg.window,
			pkg.context,
		)
	})

	ioc.RegisterDependency[logger.Logger, Api](b)
	ioc.RegisterDependency[events.Builder, Api](b)

	ioc.WrapService(b, runtimeservice.OrderCleanUp, func(c ioc.Dic, b runtimeservice.Builder) runtimeservice.Builder {
		b.OnStop(func(r runtimeservice.Runtime) {
			api := ioc.Get[Api](c)
			sdl.GLDeleteContext(api.Ctx())
			_ = api.Window().Destroy()
			sdl.Quit()
		})
		return b
	})
	ioc.RegisterDependency[runtimeservice.Builder, Api](b)
}
