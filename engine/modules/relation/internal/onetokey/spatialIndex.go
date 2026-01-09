package onetokey

import (
	"engine/modules/relation"
	"engine/services/datastructures"
	"engine/services/ecs"
)

type spatialRelation[IndexType any] struct {
	world    ecs.World
	dirtySet ecs.DirtySet

	entities datastructures.SparseArray[ecs.EntityID, uint32]
	indices  datastructures.SparseArray[uint32, ecs.EntityID]

	componentIndex func(ecs.EntityID) (IndexType, bool)
	indexNumber    func(IndexType) uint32
}

func newSpatialIndex[IndexType any](
	w ecs.World,
	dirtySet ecs.DirtySet,
	componentIndex func(ecs.EntityID) (IndexType, bool),
	indexNumber func(IndexType) uint32,
) relation.EntityToKeyTool[IndexType] {
	indexGlobal := &spatialRelation[IndexType]{
		world:    w,
		dirtySet: dirtySet,

		entities: datastructures.NewSparseArray[ecs.EntityID, uint32](),
		indices:  datastructures.NewSparseArray[uint32, ecs.EntityID](),

		componentIndex: componentIndex,
		indexNumber:    indexNumber,
	}
	w.SaveGlobal(indexGlobal)

	return indexGlobal
}

func (i *spatialRelation[IndexType]) Get(index IndexType) (ecs.EntityID, bool) {
	for _, entity := range i.dirtySet.Get() {
		indexType, ok := i.componentIndex(entity)
		if !ok {
			if number, ok := i.entities.Get(entity); ok {
				i.entities.Remove(entity)
				i.indices.Remove(number)
			}
			continue
		}
		number := i.indexNumber(indexType)
		i.entities.Set(entity, number)
		i.indices.Set(number, entity)
	}

	number := i.indexNumber(index)
	return i.indices.Get(number)
}
