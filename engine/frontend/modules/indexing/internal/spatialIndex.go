package internal

import (
	"frontend/modules/indexing"
	"shared/services/datastructures"
	"shared/services/ecs"
	"sync"
)

type spatialIndex[Component, IndexType any] struct {
	mutex    *sync.Mutex
	indices  datastructures.SparseArray[uint32, ecs.EntityID]
	entities datastructures.SparseArray[ecs.EntityID, uint32]

	componentArray ecs.ComponentsArray[Component]
	componentIndex func(Component) IndexType
	indexNumber    func(IndexType) uint32
}

func newIndex[Component, IndexType any](
	w ecs.World,
	componentIndex func(Component) IndexType,
	indexNumber func(IndexType) uint32,
) indexing.Indices[IndexType] {
	componentArray := ecs.GetComponentsArray[Component](w.Components())
	indexGlobal := spatialIndex[Component, IndexType]{
		mutex:    &sync.Mutex{},
		indices:  datastructures.NewSparseArray[uint32, ecs.EntityID](),
		entities: datastructures.NewSparseArray[ecs.EntityID, uint32](),

		componentArray: componentArray,
		componentIndex: componentIndex,
		indexNumber:    indexNumber,
	}
	w.SaveGlobal(indexGlobal)

	componentArray.OnAdd(indexGlobal.Upsert)
	componentArray.OnChange(indexGlobal.Upsert)
	componentArray.OnRemove(indexGlobal.Remove)

	return indexGlobal
}

func (i spatialIndex[Component, IndexType]) Get(index IndexType) (ecs.EntityID, bool) {
	number := i.indexNumber(index)
	return i.indices.Get(number)
}

func (i spatialIndex[Component, IndexType]) Upsert(ei []ecs.EntityID) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	for _, entity := range ei {
		component, err := i.componentArray.GetComponent(entity)
		if err != nil {
			continue
		}
		indexType := i.componentIndex(component)
		number := i.indexNumber(indexType)
		i.indices.Set(number, entity)
		i.entities.Set(entity, number)
	}
}

func (i spatialIndex[Component, IndexType]) Remove(ei []ecs.EntityID) {
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
