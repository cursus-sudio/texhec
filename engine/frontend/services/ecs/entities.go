package ecs

import "fmt"

type entity struct{}

type entitiesImpl struct {
	components *componentsImpl
	counter    int

	existingEntities map[EntityID]entity
}

func newEntities(components *componentsImpl) *entitiesImpl {
	return &entitiesImpl{
		components:       components,
		counter:          0,
		existingEntities: make(map[EntityID]entity),
	}
}

func (entitiesStorage *entitiesImpl) NewEntity() EntityID {
	// can later switch this to guid
	index := entitiesStorage.counter
	entitiesStorage.counter += 1
	id := EntityID{
		id: fmt.Sprintf("%d", index),
	}
	entitiesStorage.existingEntities[id] = entity{}
	entitiesStorage.components.AddEntity(id)
	return id
}

func (entities *entitiesImpl) RemoveEntity(entityId EntityID) {
	delete(entities.existingEntities, entityId)
	entities.components.RemoveEntity(entityId)
}

func (entitiesStorage *entitiesImpl) GetEntities() []EntityID {
	entities := make([]EntityID, 0, len(entitiesStorage.existingEntities))
	for entityId := range entitiesStorage.existingEntities {
		entities = append(entities, entityId)
	}
	return entities
}

func (entities *entitiesImpl) EntityExists(entityId EntityID) bool {
	_, ok := entities.existingEntities[entityId]
	return ok
}
