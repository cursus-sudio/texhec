package internal

import (
	"frontend/modules/indexing"
	"shared/services/ecs"
)

func NewSpatialIndexingFactory[IndexType any](
	queryFactory func(ecs.World) ecs.LiveQuery,
	componentIndexFactory func(ecs.World) func(ecs.EntityID) IndexType,
	indexNumber func(IndexType) uint32,
) ecs.ToolFactory[indexing.SpatialIndexTool[IndexType]] {
	return ecs.NewToolFactory(func(w ecs.World) indexing.SpatialIndexTool[IndexType] {
		if index, err := ecs.GetGlobal[spatialIndex[IndexType]](w); err == nil {
			return index
		}
		w.LockGlobals()
		defer w.UnlockGlobals()
		if index, err := ecs.GetGlobal[spatialIndex[IndexType]](w); err == nil {
			return index
		}
		query := queryFactory(w)
		componentIndex := componentIndexFactory(w)
		return newIndex(w, query, componentIndex, indexNumber)
	})
}
