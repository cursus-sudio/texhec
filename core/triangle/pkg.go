package triangle

import (
	"frontend/services/frames"
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
	tools, err := NewTools()
	if err != nil {
		panic(err.Error())
	}

	ioc.RegisterSingleton(b, func(c ioc.Dic) *triangleTools { return tools })

	ioc.WrapService(b, frames.Draw, func(c ioc.Dic, b events.Builder) events.Builder {
		events.Listen(b, func(e frames.FrameEvent) {
			gl.UseProgram(tools.ShaderProgram)
			gl.BindVertexArray(tools.TriangleVAO)
			// gl.DrawArrays(gl.LINE_LOOP, 0, 3)
			gl.DrawArrays(gl.TRIANGLES, 0, 3)
			gl.BindVertexArray(0)
		})
		return b
	})

	ioc.WrapService(b, appruntime.OrderCleanUp, func(c ioc.Dic, b appruntime.Builder) appruntime.Builder {
		tools := ioc.Get[*triangleTools](c)
		b.OnStop(func(r appruntime.Runtime) {
			gl.DeleteProgram(tools.ShaderProgram)
			gl.DeleteVertexArrays(1, &tools.TriangleVAO)
		})
		return b
	})
}
