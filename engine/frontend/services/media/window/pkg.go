package window

import (
	"fmt"
	"frontend/services/frames"
	"shared/services/logger"
	runtimeservice "shared/services/runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
	"github.com/veandco/go-sdl2/sdl"
)

type Pkg struct {
	window  *sdl.Window
	context sdl.GLContext
}

func Package(
	window *sdl.Window,
	context sdl.GLContext,
) Pkg {
	return Pkg{
		window:  window,
		context: context,
	}
}

var glErrorStrings = map[uint32]string{
	gl.NO_ERROR:                      "GL_NO_ERROR",
	gl.INVALID_ENUM:                  "GL_INVALID_ENUM",
	gl.INVALID_VALUE:                 "GL_INVALID_VALUE",
	gl.INVALID_OPERATION:             "GL_INVALID_OPERATION",
	gl.STACK_OVERFLOW:                "GL_STACK_OVERFLOW",
	gl.STACK_UNDERFLOW:               "GL_STACK_UNDERFLOW",
	gl.OUT_OF_MEMORY:                 "GL_OUT_OF_MEMORY",
	gl.INVALID_FRAMEBUFFER_OPERATION: "GL_INVALID_FRAMEBUFFER_OPERATION",
	gl.CONTEXT_LOST:                  "GL_CONTEXT_LOST",
	// gl.TABLE_TOO_LARGE:               "GL_TABLE_TOO_LARGE", // Less common in modern GL
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Api {
		return newApi(
			pkg.window,
			pkg.context,
		)
	})

	// TEMP
	ioc.WrapService(b, frames.Clear, func(c ioc.Dic, b frames.Builder) frames.Builder {
		return b.OnFrame(func(of frames.OnFrame) {
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		})
	})

	ioc.WrapService(b, frames.AfterDraw, func(c ioc.Dic, b frames.Builder) frames.Builder {
		logger := ioc.Get[logger.Logger](c)
		b.OnFrame(func(of frames.OnFrame) {
			if glErr := gl.GetError(); glErr != gl.NO_ERROR {
				logger.Error(fmt.Errorf("opengl error: %x %s\n", glErr, glErrorStrings[glErr]))
			}
		})
		//
		api := ioc.Get[Api](c)
		b.OnFrame(func(of frames.OnFrame) {
			api.Window().GLSwap()
		})
		return b
	})
	ioc.RegisterDependency[frames.Builder, Api](b)

	ioc.WrapService(b, runtimeservice.OrderCleanUp, func(c ioc.Dic, b runtimeservice.Builder) runtimeservice.Builder {
		b.OnStop(func(r runtimeservice.Runtime) {
			api := ioc.Get[Api](c)
			sdl.GLDeleteContext(api.Ctx())
			api.Window().Destroy()
			sdl.Quit()
		})
		return b
	})
	ioc.RegisterDependency[runtimeservice.Builder, Api](b)
}
