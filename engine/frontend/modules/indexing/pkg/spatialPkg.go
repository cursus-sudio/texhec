package indexingpkg

import (
	"frontend/modules/indexing"
	"frontend/modules/indexing/internal/indices"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type SpatialIndexTool[IndexType any] indexing.Indices[IndexType]

type spatialIndexingPkg[IndexType any] struct {
	queryFactory   func(ecs.World) ecs.LiveQuery
	componentIndex func(ecs.World) func(ecs.EntityID) (IndexType, bool)
	indexNumber    func(IndexType) uint32
}

func SpatialIndexPackage[IndexType any](
	queryFactory func(ecs.World) ecs.LiveQuery,
	componentIndex func(ecs.World) func(entity ecs.EntityID) (indexType IndexType, ok bool),
	indexNumber func(index IndexType) uint32,
) ioc.Pkg {
	return spatialIndexingPkg[IndexType]{
		queryFactory:   queryFactory,
		componentIndex: componentIndex,
		indexNumber:    indexNumber,
	}
}

func (pkg spatialIndexingPkg[IndexType]) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[indexing.Indices[IndexType]] {
		return indices.NewSpatialIndexingFactory(
			pkg.queryFactory,
			pkg.componentIndex,
			pkg.indexNumber,
		)
	})
}
