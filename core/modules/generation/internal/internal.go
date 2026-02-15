package internal

import (
	"core/modules/generation"
	"core/modules/registry"
	"core/modules/tile"
	"engine"
	"engine/modules/batcher"
	"engine/modules/collider"
	"engine/modules/grid"
	"engine/modules/inputs"
	"engine/modules/noise"
	"slices"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ogiusek/ioc/v2"
)

type config struct {
	types       []tile.Type
	tilesPerJob int
}

type service struct {
	engine.World `inject:"1"`
	GameAssets   registry.Assets `inject:"1"`
	Tile         tile.Service    `inject:"1"`

	config
}

func NewService(c ioc.Dic) generation.Service {
	s := ioc.GetServices[service](c)
	s.types = []tile.Type{}
	s.addChance(registry.TileWater, 35)
	s.addChance(registry.TileSand, 15)
	s.addChance(registry.TileGrass, 45)
	s.addChance(registry.TileMountain, 5)

	s.tilesPerJob = 100
	return &s
}

func (s *service) addChance(tileType tile.Type, chance int) {
	s.types = append(s.types, slices.Repeat([]tile.Type{tileType}, chance)...)
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
		// typesCount := []int64{
		// 	0,
		// 	0,
		// 	0,
		// 	0,
		// }
		// count := gridComponent.GetLastIndex()
		// for i := range count {
		// 	v := gridComponent.GetTile(i)
		// 	atomic.AddInt64(&typesCount[v-1], 1)
		// }
		// text := &strings.Builder{}
		// for i, typeCount := range typesCount {
		// 	fmt.Fprintf(text, "%v %05.2f%% \n",
		// 		i,
		// 		float64(typeCount*100)/float64(count),
		// 	)
		// }
		//
		// print(text.String())
		// print("\n\n\n")

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
