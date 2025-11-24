package relationpkg

import (
	"frontend/modules/relation"
	"frontend/modules/relation/internal/onetokey"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type spatialRelationPkg[IndexType any] struct {
	queryFactory   func(ecs.World) ecs.LiveQuery
	componentIndex func(ecs.World) func(ecs.EntityID) (IndexType, bool)
	indexNumber    func(IndexType) uint32
}

func SpatialRelationPackage[IndexType any](
	queryFactory func(ecs.World) ecs.LiveQuery,
	componentIndex func(ecs.World) func(entity ecs.EntityID) (indexType IndexType, ok bool),
	indexNumber func(index IndexType) uint32,
) ioc.Pkg {
	return spatialRelationPkg[IndexType]{
		queryFactory:   queryFactory,
		componentIndex: componentIndex,
		indexNumber:    indexNumber,
	}
}

func (pkg spatialRelationPkg[IndexType]) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[relation.EntityToKeyTool[IndexType]] {
		return onetokey.NewSpatialRelationFactory(
			pkg.queryFactory,
			pkg.componentIndex,
			pkg.indexNumber,
		)
	})
}
