package ecs

import "fmt"

type entity struct{}

type entitiesImpl struct {
	components *componentsImpl
	counter    int

	existingEntities map[EntityId]entity
}

func newEntities(components *componentsImpl) *entitiesImpl {
	return &entitiesImpl{
		components:       components,
		counter:          0,
		existingEntities: make(map[EntityId]entity),
	}
}

func (entitiesStorage *entitiesImpl) NewEntity() EntityId {
	// can later switch this to guid
	index := entitiesStorage.counter
	entitiesStorage.counter += 1
	id := EntityId{
		id: fmt.Sprintf("%d", index),
	}
	entitiesStorage.existingEntities[id] = entity{}
	entitiesStorage.components.entityComponents[id] = make(map[ComponentType]*Component)
	return id
}

func (entities *entitiesImpl) RemoveEntity(entityId EntityId) {
	delete(entities.existingEntities, entityId)
	delete(entities.components.entityComponents, entityId)
}

func (entitiesStorage *entitiesImpl) GetEntities() []EntityId {
	entities := make([]EntityId, 0, len(entitiesStorage.existingEntities))
	for entityId := range entitiesStorage.existingEntities {
		entities = append(entities, entityId)
	}
	return entities
}

func (entities *entitiesImpl) EntityExists(entityId EntityId) bool {
	_, ok := entities.existingEntities[entityId]
	return ok
}
