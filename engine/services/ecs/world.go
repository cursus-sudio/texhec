package ecs

import "github.com/ogiusek/events"

// interface

type World interface {
	entitiesInterface
	componentsInterface
	globalsInterface
	eventsInterface
}

//

type world struct {
	entitiesInterface
	*componentsImpl
	globalsInterface
	eventsInterface
}

func NewWorld() World {
	entitiesImpl := newEntities()
	componentsImpl := newComponents(entitiesImpl.entities)
	globalsImpl := newGlobals()
	eventsImpl := newEvents(events.NewBuilder())

	return &world{
		entitiesInterface: entitiesImpl,
		componentsImpl:    componentsImpl,
		globalsInterface:  globalsImpl,
		eventsInterface:   eventsImpl,
	}
}

func (world world) NewEntity() EntityID {
	entity := world.entitiesInterface.NewEntity()
	return entity
}

func (world world) RemoveEntity(entity EntityID) {
	world.componentsImpl.RemoveEntity(entity)
	world.entitiesInterface.RemoveEntity(entity)
}

func (world world) GetEntities() []EntityID {
	return world.entitiesInterface.GetEntities()
}

func (world world) EntityExists(entity EntityID) bool {
	return world.entitiesInterface.EntityExists(entity)
}
