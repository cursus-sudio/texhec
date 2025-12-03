package onetokey

import (
	"engine/modules/relation"
	"engine/services/ecs"
	"sync"
)

func NewMapRelationFactory[IndexType comparable](
	queryFactory func(ecs.World) ecs.LiveQuery,
	componentIndexFactory func(ecs.World) func(ecs.EntityID) (IndexType, bool),
) ecs.ToolFactory[relation.EntityToKeyTool[IndexType]] {
	mutex := &sync.Mutex{}
	return ecs.NewToolFactory(func(w ecs.World) relation.EntityToKeyTool[IndexType] {
		if index, err := ecs.GetGlobal[mapRelation[IndexType]](w); err == nil {
			return index
		}
		mutex.Lock()
		defer mutex.Unlock()
		if index, err := ecs.GetGlobal[mapRelation[IndexType]](w); err == nil {
			return index
		}
		query := queryFactory(w)
		componentIndex := componentIndexFactory(w)
		return newMapIndex(w, query, componentIndex)
	})
}
