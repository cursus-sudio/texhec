package ecs

import (
	"engine/services/datastructures"
)

// interface

type EntityID uint32

func (id EntityID) Index() int { return int(id) }

func NewEntityID(id uint64) EntityID { return EntityID(id) }

//

type entitiesInterface interface {
	GetEntities() []EntityID
	EntityExists(EntityID) bool

	NewEntity() EntityID
	EnsureExists(EntityID)
	RemoveEntity(EntityID)
}

// impl

type entitiesImpl struct {
	counter uint64
	holes   datastructures.SparseSet[EntityID]

	entities datastructures.SparseSet[EntityID]
}

func newEntities() *entitiesImpl {
	return &entitiesImpl{
		counter:  0,
		holes:    datastructures.NewSparseSet[EntityID](),
		entities: datastructures.NewSparseSet[EntityID](),
	}
}

func (entitiesStorage *entitiesImpl) GetEntities() []EntityID {
	return entitiesStorage.entities.GetIndices()
}

func (entitiesStorage *entitiesImpl) EntityExists(entity EntityID) bool {
	return entitiesStorage.entities.Get(entity)
}

func (entitiesStorage *entitiesImpl) NewEntity() EntityID {
	if holes := entitiesStorage.holes.GetIndices(); len(holes) != 0 {
		id := holes[0]
		_ = entitiesStorage.holes.Remove(0)
		entitiesStorage.entities.Add(id)
		return id
	}
	entitiesStorage.counter += 1
	index := entitiesStorage.counter
	id := EntityID(index)
	entitiesStorage.entities.Add(id)
	return id
}

func (entitiesStorage *entitiesImpl) EnsureExists(entity EntityID) {
	if ok := entitiesStorage.entities.Get(entity); ok {
		return
	}
	for i := entity; i < EntityID(entitiesStorage.counter); i++ {
		entitiesStorage.holes.Add(entity)
	}
	entitiesStorage.entities.Add(entity)
	entitiesStorage.holes.Remove(entity)
}

func (entitiesStorage *entitiesImpl) RemoveEntity(entity EntityID) {
	entitiesStorage.holes.Add(entity)
	entitiesStorage.entities.Remove(entity)
}
