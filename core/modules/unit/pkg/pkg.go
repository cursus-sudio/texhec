package unitpkg

import (
	"core/modules/unit"
	"core/modules/unit/internal"
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
	unitSize   int32
	gridDepth  float32
	gridGroups groups.GroupsComponent
}

func Package(
	unitSize int32,
	gridDepth float32,
	groups groups.GroupsComponent,
) ioc.Pkg {
	return pkg{unitSize, gridDepth, groups}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) internal.UnitRenderSystemRegister {
		return internal.NewUnitRenderSystemRegister(
			ioc.Get[texturearray.Factory](c),
			ioc.Get[logger.Logger](c),
			ioc.Get[vbo.VBOFactory[unit.UnitComponent]](c),
			ioc.Get[assets.AssetsStorage](c),
			pkg.unitSize,
			pkg.gridDepth,
			pkg.gridGroups,
			ioc.Get[ecs.ToolFactory[camera.CameraTool]](c),
		)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) unit.UnitTool {
		return ioc.Get[internal.UnitRenderSystemRegister](c)
	})
	ioc.RegisterSingleton(b, func(c ioc.Dic) unit.System {
		return ioc.Get[internal.UnitRenderSystemRegister](c)
	})

	ioc.RegisterSingleton(b, func(c ioc.Dic) vbo.VBOFactory[unit.UnitComponent] {
		return func() vbo.VBOSetter[unit.UnitComponent] {
			vbo := vbo.NewVBO[unit.UnitComponent](func() {
				var i uint32 = 0

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(unit.UnitComponent{})), uintptr(unsafe.Offsetof(unit.UnitComponent{}.Pos.X)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.INT,
					int32(unsafe.Sizeof(unit.UnitComponent{})), uintptr(unsafe.Offsetof(unit.UnitComponent{}.Pos.Y)))
				gl.EnableVertexAttribArray(i)
				i++

				gl.VertexAttribIPointerWithOffset(i, 1, gl.UNSIGNED_INT,
					int32(unsafe.Sizeof(unit.UnitComponent{})), uintptr(unsafe.Offsetof(unit.UnitComponent{}.Type)))
				gl.EnableVertexAttribArray(i)
				i++
			})
			return vbo
		}
	})

}
