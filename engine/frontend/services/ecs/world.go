package ecs

import "sync"

type world struct {
	entitiesInterface
	*componentsImpl
	registryInterface
	mutex sync.Locker
}

func NewWorld() World {
	mutex := &sync.RWMutex{}
	componentsImpl := newComponents(mutex)
	entitiesImpl := newEntities(mutex)
	registryImpl := newRegistry(&sync.RWMutex{})

	return &world{
		entitiesInterface: entitiesImpl,
		componentsImpl:    componentsImpl,
		registryInterface: registryImpl,
	}
}

func (world world) NewEntity() EntityID {
	entity := world.entitiesInterface.NewEntity()
	world.componentsImpl.AddEntity(entity)
	return entity
}

func (world world) RemoveEntity(entity EntityID) {
	world.entitiesInterface.RemoveEntity(entity)
	world.componentsImpl.RemoveEntity(entity)
}

func (world world) GetEntities() []EntityID {
	return world.entitiesInterface.GetEntities()
}

func (world world) EntityExists(entity EntityID) bool {
	return world.entitiesInterface.EntityExists(entity)
}

func (world world) SaveComponent(entity EntityID, component Component) error {
	return world.componentsImpl.SaveComponent(entity, component)
}

func (world world) GetComponent(entityId EntityID, componentType ComponentType) (Component, error) {
	return world.componentsImpl.GetComponent(entityId, componentType)
}

func (world world) RemoveComponent(entity EntityID, componentType ComponentType) {
	world.componentsImpl.RemoveComponent(entity, componentType)
}

func (world world) QueryEntitiesWithComponents(componentTypes ...ComponentType) LiveQuery {
	return world.componentsImpl.QueryEntitiesWithComponents(componentTypes...)
}
