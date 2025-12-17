package genericrendererpkg

import (
	"engine/modules/camera"
	"engine/modules/genericrenderer"
	"engine/modules/genericrenderer/internal/renderer"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/ecs"
	"engine/services/graphics/texture"
	"engine/services/graphics/vao/vbo"
	"engine/services/logger"
	"engine/services/media/window"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[genericrenderer.Vertex] {
		return func() vbo.VBOSetter[genericrenderer.Vertex] {
			vbo := vbo.NewVBO[genericrenderer.Vertex](func() {
				gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false,
					int32(unsafe.Sizeof(genericrenderer.Vertex{})), uintptr(unsafe.Offsetof(genericrenderer.Vertex{}.Pos)))
				gl.EnableVertexAttribArray(0)

				gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false,
					int32(unsafe.Sizeof(genericrenderer.Vertex{})), uintptr(unsafe.Offsetof(genericrenderer.Vertex{}.TexturePos)))
				gl.EnableVertexAttribArray(1)
			})
			return vbo
		}
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) genericrenderer.System {
		return renderer.NewSystem(
			ioc.Get[window.Api](c),
			ioc.Get[assets.AssetsStorage](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[vbo.VBOFactory[genericrenderer.Vertex]](c),
			ioc.Get[texture.Factory](c),
			ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
			ioc.Get[ecs.ToolFactory[transform.TransformTool]](c),
		)
	})
}
