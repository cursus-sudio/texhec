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

	upsertListeners []func([]ecs.EntityID)
	removeListeners []func([]ecs.EntityID)

	componentIndex func(ecs.EntityID) (IndexType, bool)
}

func newMapIndex[IndexType comparable](
	w ecs.World,
	query ecs.LiveQuery,
	componentIndex func(ecs.EntityID) (IndexType, bool),
) relation.EntityToKeyTool[IndexType] {
	indexGlobal := spatialRelation[IndexType]{
		world: w,
		mutex: &sync.Mutex{},

		upsertListeners: make([]func([]ecs.EntityID), 0),
		removeListeners: make([]func([]ecs.EntityID), 0),

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

func (i mapRelation[IndexType]) OnUpsert(listener func([]ecs.EntityID)) {
	values := make([]ecs.EntityID, 0)
	for _, entity := range i.indices {
		values = append(values, entity)
	}
	if len(values) != 0 {
		listener(values)
	}
	i.upsertListeners = append(i.upsertListeners, listener)
	i.world.SaveGlobal(i)
}

func (i mapRelation[IndexType]) OnRemove(listener func([]ecs.EntityID)) {
	i.removeListeners = append(i.removeListeners, listener)
	i.world.SaveGlobal(i)
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
	if len(added) != 0 {
		for _, listener := range i.upsertListeners {
			listener(added)
		}
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
	if len(removed) != 0 {
		for _, listener := range i.removeListeners {
			listener(removed)
		}
	}
}
