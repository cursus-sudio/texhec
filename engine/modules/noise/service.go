package noise

import (
	"engine/modules/seed"

	"github.com/go-gl/mathgl/mgl64"
)

type Noise interface {
	// returns normalized value <0, 1> with equal distribution
	Read(mgl64.Vec2) float64
}

type LayerConfig struct {
	CellSize   float64 // default size is 1
	Multiplier float64 // default value is 1
}

func NewLayer(cellSize, multiplier float64) LayerConfig {
	return LayerConfig{
		CellSize:   cellSize,
		Multiplier: multiplier,
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
