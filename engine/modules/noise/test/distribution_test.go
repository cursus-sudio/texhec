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
	}
	for testName, noiseAdd := range tests {
		t.Run(testName, func(t *testing.T) {
			noise := noiseAdd(s.Noise.NewNoise(seed.New(1)))
			s.TestDistribution(noise.Build())
		})
	}
}
