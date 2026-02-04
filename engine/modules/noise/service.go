package noise

import (
	"engine/modules/seed"

	"github.com/go-gl/mathgl/mgl64"
)

type Noise interface {
	Wrap(easing func(v float64) float64) Noise
	// returns normalized value
	Read(mgl64.Vec2) float64
}

// *-*
// | | y
// *-*
// .x
// size if for x and y value
// intensity of to normalize <0, Intensity>
type LayerConfig struct {
	CellSize        float64 // default size is 1
	ValueMultiplier float64 //
}

//

// each layer offset is `mgl64.Vec2{math.Pi, math.Pi}.Mul(i)`
type Factory interface {
	AddPerlin(...LayerConfig) Factory
	AddValue(...LayerConfig) Factory
	// NewFBM(FBMConfig) // FBM stands for Fractal Brownian Motion
	Build() Noise
}

type Service interface {
	NewNoise(seed.Seed) Factory
}
