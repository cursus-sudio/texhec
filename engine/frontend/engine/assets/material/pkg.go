package material

import (
	_ "embed"
	"frontend/services/assets"
	"frontend/services/console"
	"frontend/services/ecs"
	"frontend/services/graphics/program"
	"frontend/services/graphics/shader"
	"frontend/services/media/window"

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
	Material assets.AssetID = "materials/texturematerial"
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
				ioc.Get[console.Console](c),
				pkg.entitiesQueryAdditionalArguments,
			)
			return material.Material(), nil
		})
		return b
	})
}
