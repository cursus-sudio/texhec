package noise

import (
	"engine/modules/seed"

	"github.com/go-gl/mathgl/mgl64"
)

type Noise interface {
	// returns normalized value
	Read(mgl64.Vec2) float64
}

// *-*
// | | y
// *-*
// .x
// size if for x and y value
// intensity of to normalize <0, Intensity>
// EasingFunction is to allow custom distrubution
type LayerConfig struct {
	CellSize        float64 // default size is 1
	ValueMultiplier float64 //
	// EasingFunction  func(float64) float64 // its only used for certain noises
	Offset mgl64.Vec2 // allows to offset layers so they won't align
}

//

type Factory interface {
	AddPerlin(...LayerConfig) Factory
	AddValue(...LayerConfig) Factory
	// NewFBM(FBMConfig) // FBM stands for Fractal Brownian Motion
	Build() Noise
}

type Service interface {
	NewNoise(seed.Seed) Factory
}
