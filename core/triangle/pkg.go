package triangle

import (
	"frontend/services/frames"
	appruntime "shared/services/runtime"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

type FrontendPkg struct{}

func FrontendPackage() FrontendPkg {
	return FrontendPkg{}
}

type triangleTools struct {
	ShaderProgram uint32
	TriangleVAO   uint32
}

func NewTools() (*triangleTools, error) {
	shaderProgram, err := createShaderProgram()
	if err != nil {
		panic(err.Error())
	}
	triangleVAO := createVAO()
	if err := gl.GetError(); err != gl.NO_ERROR {
		panic(err)
	}
	return &triangleTools{
		ShaderProgram: shaderProgram,
		TriangleVAO:   triangleVAO,
	}, nil
}

func (FrontendPkg) Register(b ioc.Builder) {
	tools, err := NewTools()
	if err != nil {
		panic(err.Error())
	}

	ioc.RegisterSingleton(b, func(c ioc.Dic) *triangleTools { return tools })

	ioc.WrapService(b, frames.Draw, func(c ioc.Dic, b frames.Builder) frames.Builder {
		b.OnFrame(func(of frames.OnFrame) {
			gl.UseProgram(tools.ShaderProgram)    // Use our compiled shader program
			gl.BindVertexArray(tools.TriangleVAO) // Bind the VAO containing triangle data
			gl.DrawArrays(gl.TRIANGLES, 0, 3)     // Draw 3 vertices as a triangle
			gl.BindVertexArray(0)                 // Unbind VAO
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
