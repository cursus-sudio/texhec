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
	if layer.Easing == nil {
		layer.Easing = func(v float64) float64 { return v }
	}
	noise := fn(
		seed.Value(),
		mgl64.Vec2{math.Pi, math.Pi}.Mul(float64(i)),
		layer,
	)
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
	return noise.MergeNoise(f.Noises...)
}
