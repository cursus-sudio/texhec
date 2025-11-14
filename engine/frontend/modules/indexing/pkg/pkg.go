package indexingpkg

import (
	"frontend/modules/indexing"
	"frontend/modules/indexing/internal"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type spatialIndexingPkg[IndexType any] struct {
	queryFactory   func(ecs.World) ecs.LiveQuery
	componentIndex func(ecs.World) func(ecs.EntityID) IndexType
	indexNumber    func(IndexType) uint32
}

func SpatialIndexingPackage[IndexType any](
	queryFactory func(ecs.World) ecs.LiveQuery,
	componentIndex func(ecs.World) func(entity ecs.EntityID) IndexType,
	indexNumber func(index IndexType) uint32,
) ioc.Pkg {
	return spatialIndexingPkg[IndexType]{
		queryFactory:   queryFactory,
		componentIndex: componentIndex,
		indexNumber:    indexNumber,
	}
}

func (pkg spatialIndexingPkg[IndexType]) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[indexing.SpatialIndexTool[IndexType]] {
		return internal.NewSpatialIndexingFactory(
			pkg.queryFactory,
			pkg.componentIndex,
			pkg.indexNumber,
		)
	})
}
