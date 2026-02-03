package internal

import (
	"engine/modules/noise"
	"engine/modules/seed"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type factory struct {
	Logger logger.Logger `inject:"1"`
	Seed   seed.Seed
	Noises []Noise
}

func NewFactory(c ioc.Dic, seed seed.Seed) noise.Factory {
	f := ioc.GetServices[*factory](c)
	f.Seed = seed
	f.Noises = make([]Noise, 0)
	return f
}

func (f *factory) AddPerlin(layers ...noise.LayerConfig) noise.Factory {
	for _, layer := range layers {
		seed := seed.New(f.Seed.Value() + uint64(len(f.Noises)))
		f.Noises = append(f.Noises, NewPerlinNoise(seed.Value(), layer))
	}
	return f
}

func (f *factory) AddValue(layers ...noise.LayerConfig) noise.Factory {
	for _, layer := range layers {
		seed := seed.New(f.Seed.Value() + uint64(len(f.Noises)))
		f.Noises = append(f.Noises, NewValueNoise(seed.Value(), layer))
	}
	return f
}

func (f *factory) Build() noise.Noise {
	return MergeNoise(f.Noises...)
}
