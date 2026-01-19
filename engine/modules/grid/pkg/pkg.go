package gridpkg

import (
	"engine/modules/collider"
	"engine/modules/grid"
	"engine/modules/grid/internal/gridcollider"
	"engine/modules/grid/internal/service"
	"engine/services/codec"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type pkg[Tile grid.TileConstraint] struct {
	indexEvent func(ecs.EntityID, grid.Index) any
}

// index event can be nil if layer has no collider
func Package[Tile grid.TileConstraint](
	indexEvent func(ecs.EntityID, grid.Index) any,
) ioc.Pkg {
	return pkg[Tile]{
		indexEvent,
	}
}

func (pkg pkg[Tile]) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) grid.Service[Tile] {
		return service.NewService[Tile](c)
	})

	ioc.WrapService(b, func(c ioc.Dic, b codec.Builder) {
		b.
			// components
			Register(grid.SquareGridComponent[Tile]{})
	})

	ioc.WrapService(b, func(c ioc.Dic, collider collider.Service) {
		policy := gridcollider.NewColliderWithPolicy[Tile](
			c,
			pkg.indexEvent,
		)
		if policy != nil {
			collider.AddRayFallThroughPolicy(policy)
		}
	})
}
