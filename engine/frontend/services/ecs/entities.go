package ecs

import (
	"frontend/services/datastructures"
	"sync"
)

// type entity any

type entitiesImpl struct {
	counter uint64

	existingEntities datastructures.Set[EntityID]
	mutex            *sync.RWMutex
	// existingEntities map[EntityID]entity
	// cachedEntities []EntityID
}

func newEntities(mutex *sync.RWMutex) *entitiesImpl {
	return &entitiesImpl{
		counter:          0,
		existingEntities: datastructures.NewSet[EntityID](),
		mutex:            mutex,
		// existingEntities: make(map[EntityID]entity),
	}
}

func (entitiesStorage *entitiesImpl) NewEntity() EntityID {
	entitiesStorage.mutex.Lock()
	index := entitiesStorage.counter
	entitiesStorage.counter += 1
	id := EntityID(index)
	entitiesStorage.existingEntities.Add(id)
	return id
}

func (entitiesStorage *entitiesImpl) RemoveEntity(entityId EntityID) {
	entitiesStorage.mutex.Lock()
	entitiesStorage.existingEntities.RemoveElements(entityId)
}

func (entitiesStorage *entitiesImpl) GetEntities() []EntityID {
	entitiesStorage.mutex.RLocker().Lock()
	defer entitiesStorage.mutex.RLocker().Unlock()
	return entitiesStorage.existingEntities.Get()
}

func (entitiesStorage *entitiesImpl) EntityExists(entityId EntityID) bool {
	_, ok := entitiesStorage.existingEntities.GetIndex(entityId)
	return ok
}
