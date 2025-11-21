package tilepkg

import (
	"core/modules/tile"
	"core/modules/tile/internal/tilecollider"
	"core/modules/tile/internal/tilerenderer"
	"core/modules/tile/internal/tileui"
	"frontend/modules/collider"
	"frontend/modules/groups"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type pkg struct {
	pkgs []ioc.Pkg
}

func Package(
	tileSize int32,
	gridDepth float32,
	tileGroups groups.GroupsComponent,
	colliderComponent collider.ColliderComponent,
	mainLayer tile.Layer,
	layers []tile.Layer,
	minX, maxX, minY, maxY, minZ, maxZ int32,
) ioc.Pkg {
	return pkg{
		[]ioc.Pkg{
			tilecollider.Package(
				tileSize,
				gridDepth,
				tileGroups,
				colliderComponent,
				mainLayer,
				layers,
				minX, maxX, minY, maxY, minZ,
			),
			tilerenderer.Package(
				tileSize,
				gridDepth,
				tileGroups,
			),
			tileui.Package(),
		},
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) tile.System {
		return ecs.NewSystemRegister(func(w ecs.World) error {
			return nil
		})
	})
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
