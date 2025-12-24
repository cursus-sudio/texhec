package relationpkg

import (
	"engine/modules/relation"
	"engine/modules/relation/internal/onetokey"
	"engine/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type mapRelationPkg[IndexType comparable] struct {
	queryFactory   func(ecs.World) ecs.DirtySet
	componentIndex func(ecs.World) func(ecs.EntityID) (IndexType, bool)
}

func MapRelationPackage[IndexType comparable](
	queryFactory func(ecs.World) ecs.DirtySet,
	componentIndex func(ecs.World) func(entity ecs.EntityID) (indexType IndexType, ok bool),
) ioc.Pkg {
	return mapRelationPkg[IndexType]{
		queryFactory:   queryFactory,
		componentIndex: componentIndex,
	}
}

func (pkg mapRelationPkg[IndexType]) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) relation.ToolFactory[IndexType] {
		return onetokey.NewMapRelationFactory(
			pkg.queryFactory,
			pkg.componentIndex,
		)
	})
}
