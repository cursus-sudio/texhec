package internal

import (
	"core/modules/definitions"
	"core/modules/generation"
	"core/modules/tile"
	"engine"
	"engine/modules/batcher"
	"engine/modules/collider"
	"engine/modules/grid"
	"engine/modules/inputs"
	"engine/modules/noise"
	"engine/services/ecs"
	"fmt"
	"slices"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ogiusek/ioc/v2"
)

type config struct {
	types       []tile.ID
	tilesPerJob int
}

type service struct {
	engine.World `inject:"1"`
	GameAssets   definitions.Definitions `inject:"1"`
	Tile         tile.Service            `inject:"1"`

	config
}

func NewService(c ioc.Dic) generation.Service {
	s := ioc.GetServices[service](c)
	s.types = []tile.ID{}
	s.addChance(s.GameAssets.Tiles.Water, 35)
	s.addChance(s.GameAssets.Tiles.Sand, 15)
	s.addChance(s.GameAssets.Tiles.Grass, 45)
	s.addChance(s.GameAssets.Tiles.Mountain, 5)

	s.tilesPerJob = 100
	return &s
}

func (s *service) addChance(tileType ecs.EntityID, chance int) {
	tileComp, ok := s.Tile.Tile().Get(tileType)
	if !ok {
		s.Logger.Warn(fmt.Errorf("\"%v\" isn't a tile tile and therefor cannot be used in generation", tileType))
		return
	}
	s.types = append(s.types, slices.Repeat([]tile.ID{tileComp.ID}, chance)...)
}

func MapRange(val, min, max float64) float64 { return min + (val * (max - min)) }

func (s *service) Generate(c generation.Config) batcher.Task {
	task := s.Batcher.NewTask()

	multiplier := 1. / 4

	noise := s.Noise.NewNoise(c.Seed).AddValue(
		noise.NewLayer(100*multiplier, .10),
		noise.NewLayer(100*multiplier, .10),
		noise.NewLayer(040*multiplier, .10),
		noise.NewLayer(040*multiplier, .05),
		noise.NewLayer(040*multiplier, .05),
	).AddPerlin(
		noise.NewLayer(500*multiplier, .50),
		noise.NewLayer(500*multiplier, .50),
		noise.NewLayer(500*multiplier, .50),
		noise.NewLayer(500*multiplier, .50),
		noise.NewLayer(500*multiplier, .50),
		noise.NewLayer(500*multiplier, .50),
		noise.NewLayer(500*multiplier, .50),
		noise.NewLayer(100*multiplier, .20),
		//
		noise.NewLayer(040*multiplier, .05),
		noise.NewLayer(020*multiplier, .05),
	).Build()

	gridComponent := tile.NewGrid(c.Size.Coords())
	jobs := int(gridComponent.GetLastIndex()) / s.tilesPerJob
	generateBatch := batcher.NewBatch(jobs, func(i int) {
		for j := range s.tilesPerJob {
			gridI := grid.Index(i*s.tilesPerJob + j)
			s.SetTile(gridComponent, gridI, noise)
		}
	})
	flushBatch := batcher.NewBatch(1, func(i int) {
		size := s.Tile.GetTileSize()
		size.Size[0] *= float32(c.Size.X)
		size.Size[1] *= float32(c.Size.Y)

		s.Transform.Size().Set(c.Entity, size)

		s.Collider.Component().Set(c.Entity, collider.NewCollider(s.GameAssets.SquareCollider))
		s.Inputs.Stack().Set(c.Entity, inputs.StackComponent{})
		s.Tile.Grid().Set(c.Entity, gridComponent)
	})

	task.AddConcurrentBatch(generateBatch)
	task.AddOrderedBatch(flushBatch)

	return task.Build()
}

func (s *service) SetTile(
	grid grid.SquareGridComponent[tile.ID],
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
