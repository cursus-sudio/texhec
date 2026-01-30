package internal

import (
	"core/modules/generation"
	"core/modules/tile"
	"engine"
	"engine/modules/batcher"
	"engine/modules/grid"
	"math/rand/v2"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	engine.World `inject:"1"`
	Tile         tile.Service `inject:"1"`
}

func NewService(c ioc.Dic) generation.Service {
	s := ioc.GetServices[*service](c)
	return s
}

func (s *service) Generate(c generation.Configuration) batcher.Task {
	task := s.Batcher.NewTask()

	gridComponent := tile.NewGrid(c.Size.Coords())
	generateBatch := batcher.NewBatch(int(gridComponent.GetLastIndex()), func(i int) {
		rand := rand.New(rand.NewPCG(c.Seed, uint64(i)))
		tile := s.GetTile(rand)
		gridComponent.SetTile(grid.Index(i), tile)
	})
	flushBatch := batcher.NewBatch(1, func(i int) {
		s.Tile.Grid().Set(c.Entity, gridComponent)
	})

	task.AddConcurrentBatch(generateBatch)
	task.AddOrderedBatch(flushBatch)

	return task.Build()
}

var tiles []tile.Type = []tile.Type{
	tile.TileGrass,
	tile.TileWater,
	tile.TileMountain,
	tile.TileSand,
}

func (s *service) GetTile(rand *rand.Rand) tile.Type {
	i := rand.IntN(len(tiles))
	return tiles[i]
}
