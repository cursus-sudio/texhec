package ecs

import (
	"fmt"
	"frontend/services/datastructures"
)

// type entity any

type entitiesImpl struct {
	components *componentsImpl
	counter    int

	existingEntities datastructures.Set[EntityID]
	// existingEntities map[EntityID]entity
	// cachedEntities []EntityID
}

func newEntities(components *componentsImpl) *entitiesImpl {
	return &entitiesImpl{
		components:       components,
		counter:          0,
		existingEntities: datastructures.NewSet[EntityID](),
		// existingEntities: make(map[EntityID]entity),
	}
}

func (entitiesStorage *entitiesImpl) NewEntity() EntityID {
	// can later switch this to guid
	index := entitiesStorage.counter
	entitiesStorage.counter += 1
	id := EntityID{
		id: fmt.Sprintf("%d", index),
	}
	entitiesStorage.existingEntities.Add(id)
	// entitiesStorage.existingEntities[id] = nil
	entitiesStorage.components.AddEntity(id)
	// if entitiesStorage.cachedEntities != nil {
	// 	entitiesStorage.cachedEntities = append(entitiesStorage.cachedEntities, id)
	// }
	return id
}

func (entities *entitiesImpl) RemoveEntity(entityId EntityID) {
	if index, ok := entities.existingEntities.GetIndex(entityId); ok {
		entities.existingEntities.Remove(index)
		// i2, ok2 := entities.existingEntities.GetIndex(entityId)
		// if ok2 {
		// 	panic(fmt.Sprintf("i1 %d; i2 %d; ok2 %t;", index, i2, ok2))
		// }
		entities.components.RemoveEntity(entityId)
	}
	// delete(entities.existingEntities, entityId)
	// entities.cachedEntities = nil
}

func (entitiesStorage *entitiesImpl) GetEntities() []EntityID {
	return entitiesStorage.existingEntities.Get()
	// if entitiesStorage.cachedEntities == nil {
	// 	entities := make([]EntityID, 0, len(entitiesStorage.existingEntities))
	// 	for entityId := range entitiesStorage.existingEntities {
	// 		entities = append(entities, entityId)
	// 	}
	// 	entitiesStorage.cachedEntities = entities
	// }
	// return entitiesStorage.cachedEntities
}

func (entities *entitiesImpl) EntityExists(entityId EntityID) bool {
	_, ok := entities.existingEntities.GetIndex(entityId)
	// _, ok := entities.existingEntities[entityId]
	return ok
}
