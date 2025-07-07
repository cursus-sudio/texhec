package ecs

type world struct {
	entitiesInterface
	componentsInterface
}

func NewWorld() World {
	componentsImpl := newComponents()
	entitiesImpl := newEntities(componentsImpl)

	return &world{
		entitiesInterface:   entitiesImpl,
		componentsInterface: componentsImpl,
	}
}
