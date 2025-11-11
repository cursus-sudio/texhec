package internal

import (
	"frontend/modules/indexing"
	"shared/services/ecs"
)

func NewSpatialIndexingFactory[Component, IndexType any](
	componentIndex func(Component) IndexType,
	indexNumber func(IndexType) uint32,
) ecs.ToolFactory[indexing.SpatialIndexTool[Component, IndexType]] {
	return ecs.NewToolFactory(func(w ecs.World) indexing.SpatialIndexTool[Component, IndexType] {
		if index, err := ecs.GetGlobal[spatialIndex[Component, IndexType]](w); err == nil {
			return index
		}
		w.LockGlobals()
		defer w.UnlockGlobals()
		if index, err := ecs.GetGlobal[spatialIndex[Component, IndexType]](w); err == nil {
			return index
		}
		return newIndex(w, componentIndex, indexNumber)
	})
}
