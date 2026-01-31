package generation

import (
	"engine/modules/batcher"
	"engine/modules/grid"
	"engine/modules/seed"
	"engine/services/ecs"
)

type Config struct {
	Entity ecs.EntityID
	Seed   seed.Seed
	// will be generated <0, n)
	Size grid.Coords
}

func NewConfig(
	entity ecs.EntityID,
	seed seed.Seed,
	size grid.Coords,
) Config {
	return Config{
		entity,
		seed,
		size,
	}
}

type Service interface {
	// adds to world all grids
	Generate(Config) batcher.Task
}
