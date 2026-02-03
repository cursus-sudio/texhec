package internal

import (
	"engine/modules/noise"
	"engine/modules/seed"

	"github.com/ogiusek/ioc/v2"
)

type service struct {
	c ioc.Dic
}

func NewService(c ioc.Dic) noise.Service {
	return &service{c}
}

func (s *service) NewNoise(seed seed.Seed) noise.Factory {
	return NewFactory(s.c, seed)
}
