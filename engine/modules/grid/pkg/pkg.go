package gridpkg

import (
	"engine/modules/collider"
	"engine/modules/grid"
	"engine/modules/grid/internal/gridcollider"
	"engine/services/codec"

	"github.com/ogiusek/ioc/v2"
)

type pkg[Tile grid.TileConstraint] struct {
	indexEvent func(grid.Index) any
}

// index event can be nil if layer has no collider
func Package[Tile grid.TileConstraint](
	indexEvent func(grid.Index) any,
) ioc.Pkg {
	return pkg[Tile]{
		indexEvent,
	}
}

func (pkg pkg[Tile]) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// components
			Register(grid.SquareGridComponent[Tile]{})
	})

	if pkg.indexEvent == nil {
		ioc.WrapService(b, func(c ioc.Dic, collider collider.Service) {
			policy := gridcollider.NewColliderWithPolicy[Tile](
				c,
				pkg.indexEvent,
			)
			collider.AddRayFallThroughPolicy(policy)
		})
	}
}
