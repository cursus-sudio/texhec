package onetokey

import (
	"engine/modules/relation"
	"engine/services/ecs"
	"sync"
)

type mapRelation[IndexType comparable] struct {
	world   ecs.World
	mutex   *sync.Mutex
	indices map[IndexType]ecs.EntityID

	componentIndex func(ecs.EntityID) (IndexType, bool)
}

func newMapIndex[IndexType comparable](
	w ecs.World,
	query ecs.LiveQuery,
	componentIndex func(ecs.EntityID) (IndexType, bool),
) relation.EntityToKeyTool[IndexType] {
	indexGlobal := mapRelation[IndexType]{
		world:   w,
		mutex:   &sync.Mutex{},
		indices: make(map[IndexType]ecs.EntityID),

		componentIndex: componentIndex,
	}
	w.SaveGlobal(indexGlobal)

	query.OnAdd(indexGlobal.Upsert)
	query.OnChange(indexGlobal.Upsert)
	query.OnRemove(indexGlobal.Remove)

	return indexGlobal
}

func (i mapRelation[IndexType]) Get(index IndexType) (ecs.EntityID, bool) {
	entity, ok := i.indices[index]
	return entity, ok
}

func (i mapRelation[IndexType]) Upsert(ei []ecs.EntityID) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	added := make([]ecs.EntityID, 0, len(ei))
	for _, entity := range ei {
		indexType, ok := i.componentIndex(entity)
		if !ok {
			continue
		}
		added = append(added, entity)
		i.indices[indexType] = entity
	}
}

func (i mapRelation[IndexType]) Remove(ei []ecs.EntityID) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	removed := make([]ecs.EntityID, 0, len(ei))
	for _, entity := range ei {
		indexType, ok := i.componentIndex(entity)
		if !ok {
			continue
		}
		delete(i.indices, indexType)
		removed = append(removed, entity)
	}
}
