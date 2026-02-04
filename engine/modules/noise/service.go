package noise

import (
	"engine/modules/seed"

	"github.com/go-gl/mathgl/mgl64"
)

type Noise interface {
	// returns normalized value <0, 1> with equal distribution
	Read(mgl64.Vec2) float64
}

// *-*
// | | y
// *-*
// .x
// size if for x and y value
// intensity of to normalize <0, Intensity>
type LayerConfig struct {
	CellSize float64               // default size is 1
	Easing   func(float64) float64 //
}

func NewLayer(cellSize float64, easing func(v float64) float64) LayerConfig {
	return LayerConfig{
		CellSize: cellSize,
		Easing:   easing,
	}
}

//

// each layer offset is `mgl64.Vec2{math.Pi, math.Pi}.Mul(layerIndex)`
type Factory interface {
	AddPerlin(...LayerConfig) Factory
	AddValue(...LayerConfig) Factory
	Build() Noise
}

type Service interface {
	NewNoise(seed.Seed) Factory
}
