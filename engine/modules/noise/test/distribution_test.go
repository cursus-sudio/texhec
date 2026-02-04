package test

import (
	"engine/modules/noise"
	"engine/modules/seed"
	"testing"
)

func TestNoisesDistributions(t *testing.T) {
	s := NewSetup(t)
	tests := map[string]func(b noise.Factory) noise.Factory{
		"perlin": func(b noise.Factory) noise.Factory { return b.AddPerlin(s.Layer) },
		"value":  func(b noise.Factory) noise.Factory { return b.AddValue(s.Layer) },
		"merged <0, .5> + <0, .5>": func(b noise.Factory) noise.Factory {
			return b.
				AddPerlin(noise.NewLayer(10, .5)).
				AddPerlin(noise.NewLayer(10, .5))
		},
		"merged <0, .95> + <0, .05>": func(b noise.Factory) noise.Factory {
			return b.
				AddPerlin(noise.NewLayer(10, .95)).
				AddPerlin(noise.NewLayer(10, .05))
		},
	}
	for testName, noiseAdd := range tests {
		noise := noiseAdd(s.Noise.NewNoise(seed.New(1)))
		s.TestDistribution(testName, noise.Build())
	}
}
