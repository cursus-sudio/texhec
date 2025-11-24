package ecs

import (
	"shared/services/datastructures"
	"sync"
)

// interface

type EntityID uint32

func (id EntityID) Index() int { return int(id) }

func NewEntityID(id uint64) EntityID { return EntityID(id) }

//

type entitiesInterface interface {
	NewEntity() EntityID
	EnsureEntityExists(EntityID)
	RemoveEntity(EntityID)

	GetEntities() []EntityID
	EntityExists(EntityID) bool

	LockEntities()
	UnlockEntities()
}

// impl

type entitiesImpl struct {
	counter uint64
	holes   datastructures.Set[EntityID]
	mutex   sync.Locker

	entities datastructures.SparseSet[EntityID]
}

func newEntities() *entitiesImpl {
	return &entitiesImpl{
		counter:  0,
		holes:    datastructures.NewSet[EntityID](),
		mutex:    &sync.Mutex{},
		entities: datastructures.NewSparseSet[EntityID](),
	}
}

func (entitiesStorage *entitiesImpl) NewEntity() EntityID {
	if id, ok := entitiesStorage.holes.GetStored(0); ok {
		entitiesStorage.holes.Remove(0)
		entitiesStorage.entities.Add(id)
		return id
	}
	entitiesStorage.counter += 1
	index := entitiesStorage.counter
	id := EntityID(index)
	entitiesStorage.entities.Add(id)
	return id
}

func (entitiesStorage *entitiesImpl) EnsureEntityExists(entity EntityID) {
	if ok := entitiesStorage.entities.Get(entity); ok {
		return
	}

	for entitiesStorage.counter < uint64(entity) {
		entitiesStorage.counter++
		entitiesStorage.holes.Add(EntityID(entitiesStorage.counter))
	}

	entitiesStorage.holes.RemoveElements(entity)
	entitiesStorage.entities.Add(entity)
}

func (entitiesStorage *entitiesImpl) RemoveEntity(entity EntityID) {
	entitiesStorage.holes.Add(entity)
	entitiesStorage.entities.Remove(entity)
}

func (entitiesStorage *entitiesImpl) GetEntities() []EntityID {
	return entitiesStorage.entities.GetIndices()
}

func (entitiesStorage *entitiesImpl) EntityExists(entity EntityID) bool {
	return entitiesStorage.entities.Get(entity)
}

func (entitiesStorage *entitiesImpl) LockEntities()   { entitiesStorage.mutex.Lock() }
func (entitiesStorage *entitiesImpl) UnlockEntities() { entitiesStorage.mutex.Unlock() }
