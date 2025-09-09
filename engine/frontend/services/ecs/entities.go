package ecs

import (
	"frontend/services/datastructures"
	"sync"
)

// interface

type EntityID int

func (id EntityID) Index() int { return int(id) }

func NewEntityID(id uint64) EntityID { return EntityID(id) }

//

type entitiesInterface interface {
	NewEntity() EntityID
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

	existingEntities datastructures.Set[EntityID]
}

func newEntities() *entitiesImpl {
	return &entitiesImpl{
		counter:          0,
		holes:            datastructures.NewSet[EntityID](),
		mutex:            &sync.Mutex{},
		existingEntities: datastructures.NewSet[EntityID](),
	}
}

func (entitiesStorage *entitiesImpl) NewEntity() EntityID {
	if id, ok := entitiesStorage.holes.GetStored(0); ok {
		entitiesStorage.holes.Remove(0)
		return id
	}
	entitiesStorage.counter += 1
	index := entitiesStorage.counter
	id := EntityID(index)
	entitiesStorage.existingEntities.Add(id)
	return id
}

func (entitiesStorage *entitiesImpl) RemoveEntity(entityId EntityID) {
	entitiesStorage.holes.Add(entityId)
	entitiesStorage.existingEntities.RemoveElements(entityId)
}

func (entitiesStorage *entitiesImpl) GetEntities() []EntityID {
	return entitiesStorage.existingEntities.Get()
}

func (entitiesStorage *entitiesImpl) EntityExists(entityId EntityID) bool {
	_, ok := entitiesStorage.existingEntities.GetIndex(entityId)
	return ok
}

func (entitiesStorage *entitiesImpl) LockEntities()   { entitiesStorage.mutex.Lock() }
func (entitiesStorage *entitiesImpl) UnlockEntities() { entitiesStorage.mutex.Unlock() }
