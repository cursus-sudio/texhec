package ecs

import (
	"frontend/services/datastructures"
	"sync"
)

// type entity any

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
