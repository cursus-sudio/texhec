package indexingpkg

import (
	"frontend/modules/indexing"
	"frontend/modules/indexing/internal"
	"shared/services/ecs"

	"github.com/ogiusek/ioc/v2"
)

type spatialIndexingPkg[Component, IndexType any] struct {
	componentIndex func(Component) IndexType
	indexNumber    func(IndexType) uint32
}

func SpatialIndexingPackage[Component, IndexType any](
	componentIndex func(component Component) IndexType,
	indexNumber func(index IndexType) uint32,
) ioc.Pkg {
	return spatialIndexingPkg[Component, IndexType]{
		componentIndex: componentIndex,
		indexNumber:    indexNumber,
	}
}

func (pkg spatialIndexingPkg[Component, IndexType]) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) ecs.ToolFactory[indexing.SpatialIndexTool[Component, IndexType]] {
		return internal.NewSpatialIndexingFactory(
			pkg.componentIndex,
			pkg.indexNumber,
		)
	})
}
