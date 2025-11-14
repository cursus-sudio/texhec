package tilepkg

import (
	"core/modules/tilerenderer"
	"core/modules/tilerenderer/internal"
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
	ioc.RegisterSingleton(b, func(c ioc.Dic) internal.TileRenderSystemRegister {
		return internal.NewTileRenderSystemRegister(
			ioc.Get[texturearray.Factory](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[vbo.VBOFactory[internal.TileData]](c),
			ioc.Get[assets.AssetsStorage](c),
			pkg.tileSize,
			pkg.gridDepth,
			pkg.gridGroups,
			ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) tilerenderer.TileTool {
		return ioc.Get[internal.TileRenderSystemRegister](c)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) tilerenderer.System {
		return ioc.Get[internal.TileRenderSystemRegister](c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[internal.TileData] {
		return func() vbo.VBOSetter[internal.TileData] {
			vbo := vbo.NewVBO[internal.TileData](func() {
				var i uint32 = 0

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(internal.TileData{})), uintptr(unsafe.Offsetof(internal.TileData{}.Pos.X)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(internal.TileData{})), uintptr(unsafe.Offsetof(internal.TileData{}.Pos.Y)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(internal.TileData{})), uintptr(unsafe.Offsetof(internal.TileData{}.Pos.Z)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.UNSIGNED_INT,
					int32(unsafe.Sizeof(internal.TileData{})), uintptr(unsafe.Offsetof(internal.TileData{}.Type)))
				gl.EnableVertexAttribArray(i)
				i++
			})
			return vbo
		}
	})

}
