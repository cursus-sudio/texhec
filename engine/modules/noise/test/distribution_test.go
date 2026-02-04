package test

import (
	"engine/modules/seed"
	"testing"
)

func TestPerlinNoiseDistribution(t *testing.T) {
	s := NewSetup(t)
	noise := s.Noise.NewNoise(seed.New(1)).
		AddPerlin(s.Layer).
		Build()
	s.TestDistribution(noise)
}

func TestValueNoiseDistribution(t *testing.T) {
	s := NewSetup(t)
	noise := s.Noise.NewNoise(seed.New(1)).
		AddValue(s.Layer).
		Build()
	s.TestDistribution(noise)
}
