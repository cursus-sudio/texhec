package internal

import (
	"core/modules/generation"
	"core/modules/tile"
	"engine"
	"engine/modules/batcher"
	"engine/modules/grid"
	"fmt"
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
		var tileType tile.Type
		switch rand.IntN(3) {
		case 0:
			tileType = tile.TileGrass
		case 1:
			tileType = tile.TileWater
		case 2:
			tileType = tile.TileSand
		default:
			s.Logger.Warn(fmt.Errorf("missing number handler"))
		}
		gridComponent.SetTile(grid.Index(i), tile.Type(tileType))
	})
	flushBatch := batcher.NewBatch(1, func(i int) {
		s.Tile.Grid().Set(c.Entity, gridComponent)
	})

	task.AddConcurrentBatch(generateBatch)
	task.AddOrderedBatch(flushBatch)

	return task.Build()
}
