package tilerenderer

import (
	"core/modules/tile"
	"engine/services/graphics/vao/vbo"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
}

// TODO
// currently doesn't support animated tiles
// always renderes first frame if something is animated
func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) *TileRenderSystemRegister {
		return NewTileRenderSystemRegister(c)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) tile.TileAssets {
		return ioc.Get[*TileRenderSystemRegister](c)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) tile.SystemRenderer {
		return ioc.Get[*TileRenderSystemRegister](c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[tile.Type] {
		return func() vbo.VBOSetter[tile.Type] {
			vbo := vbo.NewVBO[tile.Type](func() {
				var i uint32 = 0

				gl.VertexAttribIPointerWithOffset(i, 1, gl.UNSIGNED_BYTE,
					int32(unsafe.Sizeof(tile.Type(0))), uintptr(0))
				gl.EnableVertexAttribArray(i)
			})
			return vbo
		}
	})
}
