package ecs

type world struct {
	entitiesInterface
	componentsInterface
	systemsInterface
}

func newWorld() World {
	componentsImpl := newComponents()
	entitiesImpl := newEntities(componentsImpl)
	systemsImpl := newSystems()

	return &world{
		entitiesInterface:   entitiesImpl,
		componentsInterface: componentsImpl,
		systemsInterface:    systemsImpl,
	}
}
