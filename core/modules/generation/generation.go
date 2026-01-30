package generation

import (
	"engine/modules/batcher"
	"engine/modules/grid"
	"engine/services/ecs"
)

type Configuration struct {
	Entity ecs.EntityID
	Seed   uint64
	// will be generated <0, n)
	Size grid.Coords
}

func NewConfiguration(
	entity ecs.EntityID,
	seed uint64,
	size grid.Coords,
) Configuration {
	return Configuration{
		entity,
		seed,
		size,
	}
}

type Service interface {
	// adds to world all grids
	Generate(Configuration) batcher.Task
}
