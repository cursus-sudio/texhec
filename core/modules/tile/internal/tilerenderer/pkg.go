package tilerenderer

import (
	"core/modules/tile"
	"frontend/modules/camera"
	"frontend/modules/groups"
	"frontend/services/assets"
	"frontend/services/graphics/texturearray"
	"frontend/services/graphics/vao/vbo"
	"shared/services/ecs"
	"shared/services/logger"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	tileSize   int32
	gridDepth  float32
	gridGroups groups.GroupsComponent
}

// TODO
// currently doesn't support animated tiles
// always renderes first frame if something is animated
func Package(
	tileSize int32,
	gridDepth float32,
	groups groups.GroupsComponent,
) ioc.Pkg {
	return pkg{tileSize, gridDepth, groups}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) TileRenderSystemRegister {
		return NewTileRenderSystemRegister(
			ioc.Get[texturearray.Factory](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[vbo.VBOFactory[TileData]](c),
			ioc.Get[assets.AssetsStorage](c),
			pkg.tileSize,
			pkg.gridDepth,
			pkg.gridGroups,
			ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) tile.TileAssets {
		return ioc.Get[TileRenderSystemRegister](c)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) tile.SystemRenderer {
		return ioc.Get[TileRenderSystemRegister](c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[TileData] {
		return func() vbo.VBOSetter[TileData] {
			vbo := vbo.NewVBO[TileData](func() {
				var i uint32 = 0

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(TileData{})), uintptr(unsafe.Offsetof(TileData{}.Pos.X)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(TileData{})), uintptr(unsafe.Offsetof(TileData{}.Pos.Y)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(TileData{})), uintptr(unsafe.Offsetof(TileData{}.Pos.Z)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.UNSIGNED_INT,
					int32(unsafe.Sizeof(TileData{})), uintptr(unsafe.Offsetof(TileData{}.Type)))
				gl.EnableVertexAttribArray(i)
				i++
			})
			return vbo
		}
	})
}
