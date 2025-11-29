package onetokey

import (
	"engine/modules/relation"
	"engine/services/ecs"
	"sync"
)

func NewSpatialRelationFactory[IndexType any](
	queryFactory func(ecs.World) ecs.LiveQuery,
	componentIndexFactory func(ecs.World) func(ecs.EntityID) (IndexType, bool),
	indexNumber func(IndexType) uint32,
) ecs.ToolFactory[relation.EntityToKeyTool[IndexType]] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) relation.EntityToKeyTool[IndexType] {
		if index, err := ecs.GetGlobal[spatialRelation[IndexType]](w); err == nil {
			return index
		}
		mutex.Lock()
		defer mutex.Unlock()
		if index, err := ecs.GetGlobal[spatialRelation[IndexType]](w); err == nil {
			return index
		}
		query := queryFactory(w)
		componentIndex := componentIndexFactory(w)
		return newIndex(w, query, componentIndex, indexNumber)
	})
}
