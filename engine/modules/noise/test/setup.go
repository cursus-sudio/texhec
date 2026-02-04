package test

import (
	"engine/modules/noise"
	noisepkg "engine/modules/noise/pkg"
	"engine/services/clock"
	"engine/services/logger"
	"math"
	"testing"
	"time"

	"github.com/go-gl/mathgl/mgl64"
	"github.com/ogiusek/ioc/v2"
)

type Setup struct {
	Noise noise.Service `inject:"1"`
	T     *testing.T
	Layer noise.LayerConfig
}

func NewSetup(t *testing.T) Setup {
	b := ioc.NewBuilder()
	for _, pkg := range []ioc.Pkg{
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
		clock.Package(time.RFC3339Nano),
		noisepkg.Package(),
	} {
		pkg.Register(b)
	}
	c := b.Build()
	setup := ioc.GetServices[Setup](c)
	setup.T = t
	setup.Layer = noise.LayerConfig{
		CellSize:        10.,
		ValueMultiplier: 1,
	}
	return setup
}

func (s Setup) CalculateDistribution(
	noise noise.Noise,
	samplesSqrt int,
) [3]float64 {
	res := [3]float64{}
	for x := range samplesSqrt {
		for y := range samplesSqrt {
			v := noise.Read(mgl64.Vec2{float64(x), float64(y)})
			res[min(int(v*3), 2)]++
		}
	}

	for i := range 3 {
		res[i] /= float64(samplesSqrt * samplesSqrt)
	}

	return res
}

func (s Setup) TestDistribution(noise noise.Noise) {
	s.T.Helper()
	distribution := s.CalculateDistribution(noise, 1000)
	for _, value := range distribution {
		const target = 1. / 3
		const tolerance = .02
		if math.Abs(target-value) > tolerance {
			s.T.Errorf(
				"expected distribution to have only values %.2f ± %.2f but got %v",
				target,
				tolerance,
				distribution,
			)
			return
		}
	}
}
