package ecs

// interface

type World interface {
	entitiesInterface
	componentsInterface
	registryInterface
}

//

type world struct {
	entitiesInterface
	*componentsImpl
	registryInterface
}

func NewWorld() World {
	entitiesImpl := newEntities()
	componentsImpl := newComponents(entitiesImpl.GetEntities)
	registryImpl := newRegistry()

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

//
