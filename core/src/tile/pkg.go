package tile

import (
	"frontend/engine/components/groups"
	"frontend/engine/tools/cameras"
	"frontend/services/assets"
	"frontend/services/graphics/texturearray"
	"frontend/services/graphics/vao/vbo"
	"shared/services/ecs"
	"shared/services/logger"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

type Pkg struct {
	tileSize   int32
	gridDepth  float32
	gridGroups groups.Groups
}

func Package(tileSize int32, gridDepth float32, groups groups.Groups) ioc.Pkg {
	return Pkg{tileSize, gridDepth, groups}
}

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.RegisterTransient(b, func(c ioc.Dic) TileRenderSystemRegister {
		return newTileRenderSystemRegister(
			ioc.Get[texturearray.Factory](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[vbo.VBOFactory[TileComponent]](c),
			ioc.Get[assets.AssetsStorage](c),
			pkg.tileSize,
			pkg.gridDepth,
			pkg.gridGroups,
			ioc.Get[ecs.ToolFactory[cameras.CameraResolver]](c),
		)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[TileComponent] {
		return func() vbo.VBOSetter[TileComponent] {
			vbo := vbo.NewVBO[TileComponent](func() {
				var i uint32 = 0

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(TileComponent{})), uintptr(unsafe.Offsetof(TileComponent{}.Pos.X)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(TileComponent{})), uintptr(unsafe.Offsetof(TileComponent{}.Pos.Y)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.UNSIGNED_INT,
					int32(unsafe.Sizeof(TileComponent{})), uintptr(unsafe.Offsetof(TileComponent{}.Type)))
				gl.EnableVertexAttribArray(i)
				i++
			})
			return vbo
		}
	})

}
