package triangle

import (
	"frontend/services/frames"
	"frontend/services/media/window"
	appruntime "shared/services/runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/events"
	"github.com/ogiusek/ioc/v2"
)

type FrontendPkg struct{}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

func (FrontendPkg) Register(b ioc.Builder) {
	tools, err := NewTriangleTools()
	if err != nil {
		panic(err.Error())
	}

	ioc.RegisterSingleton(b, func(c ioc.Dic) *triangleTools { return tools })

	ioc.WrapService(b, frames.Draw, func(c ioc.Dic, b events.Builder) events.Builder {
		window := ioc.Get[window.Api](c).Window()
		resolutionLocation := gl.GetUniformLocation(tools.Program.ID, gl.Str("resolution\x00"))

		events.Listen(b, func(e frames.FrameEvent) {
			tools.Program.Draw(func() {
				tools.Texture.Draw(func() {
					width, height := window.GetSize()
					gl.Uniform2f(resolutionLocation, float32(width), float32(height))

					tools.Program.Draw(tools.VAO.Draw)
				})
			})
		})
		return b
	})

	ioc.WrapService(b, appruntime.OrderCleanUp, func(c ioc.Dic, b appruntime.Builder) appruntime.Builder {
		tools := ioc.Get[*triangleTools](c)
		b.OnStop(func(r appruntime.Runtime) {
			tools.Program.Release()
			tools.VAO.Release()
		})
		return b
	})
}
