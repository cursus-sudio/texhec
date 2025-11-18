package tilepkg

import (
	"core/modules/tile"
	"core/modules/tile/internal/tilecollider"
	"core/modules/tile/internal/tilerenderer"
	"frontend/modules/collider"
	"frontend/modules/groups"
	"shared/services/datastructures"

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
	layerEvents datastructures.SparseArray[tile.Layer, []any],
	minX, maxX, minY, maxY, minZ int32,
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
				layerEvents,
				minX, maxX, minY, maxY, minZ,
			),
			tilerenderer.Package(
				tileSize,
				gridDepth,
				tileGroups,
			),
		},
	}
}

func (pkg pkg) Register(b ioc.Builder) {
	for _, pkg := range pkg.pkgs {
		pkg.Register(b)
	}
}
