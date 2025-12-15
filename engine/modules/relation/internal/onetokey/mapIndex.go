package onetokey

import (
	"engine/modules/relation"
	"engine/services/datastructures"
	"engine/services/ecs"
)

type mapRelation[IndexType comparable] struct {
	world    ecs.World
	dirtySet ecs.DirtySet

	entities datastructures.SparseArray[ecs.EntityID, IndexType]
	indices  map[IndexType]ecs.EntityID

	componentIndex func(ecs.EntityID) (IndexType, bool)
}

func newMapIndex[IndexType comparable](
	w ecs.World,
	dirtySet ecs.DirtySet,
	componentIndex func(ecs.EntityID) (IndexType, bool),
) relation.EntityToKeyTool[IndexType] {
	indexGlobal := mapRelation[IndexType]{
		world:    w,
		dirtySet: dirtySet,

		entities: datastructures.NewSparseArray[ecs.EntityID, IndexType](),
		indices:  make(map[IndexType]ecs.EntityID),

		componentIndex: componentIndex,
	}
	w.SaveGlobal(indexGlobal)

	return indexGlobal
}

func (i mapRelation[IndexType]) Get(index IndexType) (ecs.EntityID, bool) {
	for _, entity := range i.dirtySet.Get() {
		indexType, ok := i.componentIndex(entity)
		if !ok {
			if indexType, ok := i.entities.Get(entity); ok {
				i.entities.Remove(entity)
				delete(i.indices, indexType)
			}
			continue
		}

		i.entities.Set(entity, indexType)
		i.indices[indexType] = entity
	}
	entity, ok := i.indices[index]
	return entity, ok
}

func (i mapRelation[IndexType]) Upsert(ei []ecs.EntityID) {
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
