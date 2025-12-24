package onetokey

import (
	"engine/modules/relation"
	"engine/services/ecs"
	"sync"
)

func NewSpatialRelationFactory[IndexType any](
	dirtySetFactory func(ecs.World) ecs.DirtySet,
	componentIndexFactory func(ecs.World) func(ecs.EntityID) (IndexType, bool),
	indexNumber func(IndexType) uint32,
) relation.ToolFactory[IndexType] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w relation.World) relation.EntityToKeyTool[IndexType] {
		mutex.Lock()
		defer mutex.Unlock()
		if index, ok := ecs.GetGlobal[spatialRelation[IndexType]](w); ok {
			return index
		}
		dirtySet := dirtySetFactory(w)
		componentIndex := componentIndexFactory(w)
		return newSpatialIndex(w, dirtySet, componentIndex, indexNumber)
	})
}
