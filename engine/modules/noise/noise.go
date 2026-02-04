package noise

import (
	"github.com/go-gl/mathgl/mgl64"
)

type noise func(mgl64.Vec2) float64

func (n *noise) Read(coords mgl64.Vec2) float64 { return (*n)(coords) }
func (n *noise) Wrap(easing func(float64) float64) Noise {
	fn := *n
	*n = func(v mgl64.Vec2) float64 { return easing(fn(v)) }
	return n
}

//

func NewNoise(fn func(mgl64.Vec2) float64) Noise {
	n := noise(fn)
	return &n
}
