package mainpipeline

import (
	"frontend/engine/tools/worldmesh"
	"frontend/services/assets"
	"frontend/services/graphics/vao/vbo"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[Vertex] {
		return func() vbo.VBOSetter[Vertex] {
			vbo := vbo.NewVBO[Vertex](func() {
				gl.VertexAttribPointerWithOffset(0, 3, gl.FLOAT, false,
					int32(unsafe.Sizeof(Vertex{})), uintptr(unsafe.Offsetof(Vertex{}.Pos)))
				gl.EnableVertexAttribArray(0)

				gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false,
					int32(unsafe.Sizeof(Vertex{})), uintptr(unsafe.Offsetof(Vertex{}.TexturePos)))
				gl.EnableVertexAttribArray(1)
			})
			return vbo
		}
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) worldmesh.RegisterFactory[Vertex] {
		return worldmesh.NewRegisterFactory(
			ioc.Get[assets.AssetsStorage](c),
			ioc.Get[vbo.VBOFactory[Vertex]](c),
		)
	})
}
