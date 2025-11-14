package internal

import (
	"frontend/modules/indexing"
	"shared/services/datastructures"
	"shared/services/ecs"
	"sync"
)

type spatialIndex[IndexType any] struct {
	mutex    *sync.Mutex
	indices  datastructures.SparseArray[uint32, ecs.EntityID]
	entities datastructures.SparseArray[ecs.EntityID, uint32]

	componentIndex func(ecs.EntityID) IndexType
	indexNumber    func(IndexType) uint32
}

func newIndex[IndexType any](
	w ecs.World,
	query ecs.LiveQuery,
	componentIndex func(ecs.EntityID) IndexType,
	indexNumber func(IndexType) uint32,
) indexing.Indices[IndexType] {
	indexGlobal := spatialIndex[IndexType]{
		mutex:    &sync.Mutex{},
		indices:  datastructures.NewSparseArray[uint32, ecs.EntityID](),
		entities: datastructures.NewSparseArray[ecs.EntityID, uint32](),

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

func (i spatialIndex[IndexType]) Upsert(ei []ecs.EntityID) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	for _, entity := range ei {
		indexType := i.componentIndex(entity)
		number := i.indexNumber(indexType)
		i.indices.Set(number, entity)
		i.entities.Set(entity, number)
	}
}

func (i spatialIndex[IndexType]) Remove(ei []ecs.EntityID) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	for _, entity := range ei {
		number, ok := i.entities.Get(entity)
		if !ok {
			continue
		}
		i.indices.Remove(number)
		i.entities.Remove(entity)
	}
}
