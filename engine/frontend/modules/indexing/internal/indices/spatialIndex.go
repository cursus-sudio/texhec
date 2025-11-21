package indices

import (
	"frontend/modules/indexing"
	"shared/services/datastructures"
	"shared/services/ecs"
	"sync"
)

type spatialIndex[IndexType any] struct {
	world    ecs.World
	mutex    *sync.Mutex
	indices  datastructures.SparseArray[uint32, ecs.EntityID]
	entities datastructures.SparseArray[ecs.EntityID, uint32]

	upsertListeners []func([]ecs.EntityID)
	removeListeners []func([]ecs.EntityID)

	componentIndex func(ecs.EntityID) (IndexType, bool)
	indexNumber    func(IndexType) uint32
}

func newIndex[IndexType any](
	w ecs.World,
	query ecs.LiveQuery,
	componentIndex func(ecs.EntityID) (IndexType, bool),
	indexNumber func(IndexType) uint32,
) indexing.Indices[IndexType] {
	indexGlobal := spatialIndex[IndexType]{
		world:    w,
		mutex:    &sync.Mutex{},
		indices:  datastructures.NewSparseArray[uint32, ecs.EntityID](),
		entities: datastructures.NewSparseArray[ecs.EntityID, uint32](),

		upsertListeners: make([]func([]ecs.EntityID), 0),

		componentIndex: componentIndex,
		indexNumber:    indexNumber,
	}
	w.SaveGlobal(indexGlobal)

	query.OnAdd(indexGlobal.Upsert)
	query.OnChange(indexGlobal.Upsert)
	query.OnRemove(indexGlobal.Remove)

	return indexGlobal
}

func (i spatialIndex[IndexType]) Get(index IndexType) (ecs.EntityID, bool) {
	number := i.indexNumber(index)
	return i.indices.Get(number)
}

func (i spatialIndex[IndexType]) OnUpsert(listener func([]ecs.EntityID)) {
	if values := i.indices.GetValues(); len(values) != 0 {
		listener(values)
	}
	i.upsertListeners = append(i.upsertListeners, listener)
	i.world.SaveGlobal(i)
}

func (i spatialIndex[IndexType]) OnRemove(listener func([]ecs.EntityID)) {
	i.removeListeners = append(i.removeListeners, listener)
	i.world.SaveGlobal(i)
}

func (i spatialIndex[IndexType]) Upsert(ei []ecs.EntityID) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	added := make([]ecs.EntityID, 0, len(ei))
	for _, entity := range ei {
		indexType, ok := i.componentIndex(entity)
		if !ok {
			continue
		}
		added = append(added, entity)
		number := i.indexNumber(indexType)
		i.indices.Set(number, entity)
		i.entities.Set(entity, number)
	}
	if len(added) != 0 {
		for _, listener := range i.upsertListeners {
			listener(added)
		}
	}
}

func (i spatialIndex[IndexType]) Remove(ei []ecs.EntityID) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	removed := make([]ecs.EntityID, 0, len(ei))
	for _, entity := range ei {
		number, ok := i.entities.Get(entity)
		if !ok {
			continue
		}
		i.indices.Remove(number)
		i.entities.Remove(entity)
		removed = append(removed, entity)
	}
	if len(removed) != 0 {
		for _, listener := range i.removeListeners {
			listener(removed)
		}
	}
}
