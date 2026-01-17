package gridpkg

import (
	"engine/modules/collider"
	"engine/modules/grid"
	"engine/modules/grid/internal"
	"engine/services/codec"

	"github.com/ogiusek/ioc/v2"
	"golang.org/x/exp/constraints"
)

type pkg[TileType constraints.Unsigned] struct {
	indexEvent func(grid.Index) any
}

// index event can be nil if layer has no collider
func Package[TileType constraints.Unsigned](
	indexEvent func(grid.Index) any,
) ioc.Pkg {
	return pkg[TileType]{
		indexEvent,
	}
}

func (pkg pkg[TileType]) Register(b ioc.Builder) {
	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// components
			Register(grid.SquareGridComponent[TileType]{})
	})

	if pkg.indexEvent == nil {
		ioc.WrapService(b, func(c ioc.Dic, collider collider.Service) {
			policy := internal.NewColliderWithPolicy[TileType](
				c,
				pkg.indexEvent,
			)
			collider.AddRayFallThroughPolicy(policy)
		})
	}
}
