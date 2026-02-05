package internal

import (
	"engine/modules/noise"
	"engine/modules/seed"
	"engine/services/logger"
	"math"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ogiusek/ioc/v2"
)

type factory struct {
	Logger logger.Logger `inject:"1"`
	Seed   seed.Seed
	Layers []noise.LayerConfig
	Noises []noise.Noise
}

func NewFactory(c ioc.Dic, seed seed.Seed) noise.Factory {
	f := ioc.GetServices[*factory](c)
	f.Seed = seed
	f.Noises = make([]noise.Noise, 0)
	return f
}

func (f *factory) Add(
	fn func(seed uint64, offset mgl64.Vec2, layer noise.LayerConfig) noise.Noise,
	layer noise.LayerConfig,
) {
	i := len(f.Noises)
	seed := seed.New(f.Seed.Value() + uint64(i))
	noise := fn(
		seed.Value(),
		mgl64.Vec2{math.Pi, math.Pi}.Mul(float64(i)),
		layer,
	)
	f.Layers = append(f.Layers, layer)
	f.Noises = append(f.Noises, noise)
}

func (f *factory) AddPerlin(layers ...noise.LayerConfig) noise.Factory {
	for _, layer := range layers {
		f.Add(NewPerlinNoise, layer)
	}
	return f
}

func (f *factory) AddValue(layers ...noise.LayerConfig) noise.Factory {
	for _, layer := range layers {
		f.Add(NewValueNoise, layer)
	}
	return f
}

func (f *factory) Build() noise.Noise {
	var totalWeight float64
	for _, layer := range f.Layers {
		totalWeight += layer.Weight
	}
	multiplier := 1 / totalWeight

	// manual way of calculating standard deviation
	var totalVariance float64
	for _, layer := range f.Layers {
		// each noise variance is uniform so we use const
		const noiseVariance = 1. / 12.
		value := layer.Weight * multiplier
		totalVariance += value * value * noiseVariance
	}
	standardDeviation := math.Sqrt(totalVariance)
	return noise.NewNoise(func(v mgl64.Vec2) float64 {
		var s float64
		for _, noise := range f.Noises {
			s += noise.Read(v)
		}
		s *= multiplier
		s = cdf(s, standardDeviation)
		return s
	})
}
