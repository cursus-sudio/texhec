package internal

import (
	"core/modules/generation"
	"core/modules/tile"
	"engine"
	"engine/modules/batcher"
	"engine/modules/grid"
	"engine/modules/noise"
	"slices"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ogiusek/ioc/v2"
)

type config struct {
	types       []tile.Type
	tilesPerJob int64
}

type service struct {
	engine.World `inject:"1"`
	Tile         tile.Service `inject:"1"`

	config
}

func NewService(c ioc.Dic) generation.Service {
	s := ioc.GetServices[service](c)
	s.types = []tile.Type{}
	// s.addChance(tile.TileWater, 10)
	// s.addChance(tile.TileSand, 3)
	// s.addChance(tile.TileGrass, 10)
	// s.addChance(tile.TileMountain, 1)
	s.addChance(tile.TileWater, 40)
	s.addChance(tile.TileSand, 3)
	s.addChance(tile.TileGrass, 20)
	s.addChance(tile.TileMountain, 7)
	s.tilesPerJob = 100
	return &s
}

func (s *service) addChance(tileType tile.Type, chance int) {
	s.types = append(s.types, slices.Repeat([]tile.Type{tileType}, chance)...)
}

func MapRange(val, min, max float64) float64 { return min + (val * (max - min)) }

func (s *service) Generate(c generation.Config) batcher.Task {
	c.Size.X = 1000
	c.Size.Y = 1000
	task := s.Batcher.NewTask()

	noise := s.Noise.NewNoise(c.Seed).
		AddValue(noise.LayerConfig{
			CellSize:        100.,
			ValueMultiplier: .5,
		}).
		AddValue(noise.LayerConfig{
			CellSize:        30.,
			ValueMultiplier: .3,
		}).
		AddValue(noise.LayerConfig{
			CellSize:        10,
			ValueMultiplier: .1,
		}).
		AddValue(noise.LayerConfig{
			CellSize:        5,
			ValueMultiplier: .1,
		}).
		Build()

	gridComponent := tile.NewGrid(c.Size.Coords())
	jobs := int64(gridComponent.GetLastIndex()) / s.tilesPerJob
	generateBatch := batcher.NewBatch(jobs, func(i int64) {
		for j := range s.tilesPerJob {
			gridI := grid.Index(i*s.tilesPerJob + j)
			s.SetTile(gridComponent, gridI, noise)
		}
	})
	flushBatch := batcher.NewBatch(1, func(i int64) {
		s.Tile.Grid().Set(c.Entity, gridComponent)
	})

	task.AddConcurrentBatch(generateBatch)
	task.AddOrderedBatch(flushBatch)

	return task.Build()
}

func (s *service) SetTile(
	grid grid.SquareGridComponent[tile.Type],
	index grid.Index,
	noise noise.Noise,
) {
	coords := grid.GetCoords(index)
	value := noise.Read(mgl64.Vec2{float64(coords.X), float64(coords.Y)})
	value *= float64(len(s.types))
	value = min(value, float64(len(s.types)-1))
	tileValue := s.types[int(value)]
	grid.SetTile(index, tileValue)
}
