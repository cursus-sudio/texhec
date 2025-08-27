package material

import (
	_ "embed"
	"frontend/engine/tools/worldmesh"
	"frontend/services/assets"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/graphics/shader"
	"frontend/services/graphics/vao/vbo"
	"frontend/services/media/window"
	"shared/services/logger"
	"unsafe"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/ogiusek/ioc/v2"
)

//go:embed s.vert
var vertSource string

//go:embed s.frag
var fragSource string

type Pkg struct {
	entitiesQueryAdditionalArguments []ecs.ComponentType
}

func Package(
	entitiesQueryAdditionalArguments []ecs.ComponentType,
) Pkg {
	return Pkg{
		entitiesQueryAdditionalArguments: entitiesQueryAdditionalArguments,
	}
}

var (
	Material assets.AssetID = "materials/basic_material"
)

func (pkg Pkg) Register(b ioc.Builder) {
	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b assets.AssetsStorageBuilder) assets.AssetsStorageBuilder {
		b.RegisterAsset(Material, func() (any, error) {
			material := newTextureMaterial(
				func() (program.Program, error) {
					vert, err := shader.NewShader(vertSource, shader.VertexShader)
					if err != nil {
						return nil, err
					}
					frag, err := shader.NewShader(fragSource, shader.FragmentShader)
					if err != nil {
						return nil, err
					}
					p, err := program.NewProgram(vert, frag, nil)
					if err != nil {
						vert.Release()
						frag.Release()
						return nil, err
					}
					vert.Release()
					frag.Release()
					return p, nil
				},
				ioc.Get[window.Api](c),
				ioc.Get[assets.AssetsStorage](c),
				ioc.Get[logger.Logger](c),
				pkg.entitiesQueryAdditionalArguments,
			)
			return material.Material(), nil
		})
		return b
	})

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
