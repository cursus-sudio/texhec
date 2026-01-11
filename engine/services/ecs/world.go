package ecs

// interface

type World interface {
	entitiesInterface
	componentsInterface
}

//

type world struct {
	entitiesInterface
	*componentsImpl
}

func NewWorld() World {
	entitiesImpl := newEntities()
	componentsImpl := newComponents(entitiesImpl.entities)

	return &world{
		entitiesInterface: entitiesImpl,
		componentsImpl:    componentsImpl,
	}
}

func (world world) NewEntity() EntityID {
	entity := world.entitiesInterface.NewEntity()
	return entity
}

func (world world) RemoveEntity(entity EntityID) {
	world.entitiesInterface.RemoveEntity(entity)

	for _, arr := range world.storage.arrays {
		arr.Remove(entity)
	}
}

func (world world) GetEntities() []EntityID {
	return world.entitiesInterface.GetEntities()
}

func (world world) EntityExists(entity EntityID) bool {
	return world.entitiesInterface.EntityExists(entity)
}
