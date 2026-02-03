package internal

import (
	"engine/modules/noise"

	"github.com/go-gl/mathgl/mgl64"
)

type Noise func(mgl64.Vec2) float64

func NewNoise(fn func(mgl64.Vec2) float64) Noise {
	return Noise(fn)
}

func MergeNoise(noises ...Noise) noise.Noise {
	return Noise(func(c mgl64.Vec2) float64 {
		var s float64
		for _, noise := range noises {
			s += noise.read(c)
		}
		return s
	})
}

func (n Noise) read(coords mgl64.Vec2) float64 {
	return n(coords)
}

func (n Noise) Read(coords mgl64.Vec2) float64 {
	return n(coords)
}
