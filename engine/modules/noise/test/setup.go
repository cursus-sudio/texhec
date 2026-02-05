package test

import (
	"engine/modules/noise"
	"engine/modules/noise/internal"
	noisepkg "engine/modules/noise/pkg"
	"engine/services/clock"
	"engine/services/logger"
	"math"
	"testing"
	"time"

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
	setup.Layer = noise.NewLayer(1, 1)
	return setup
}

// in case tests fail you can improve algorithm using this algorithm:
// 1. Sample 10,000 Perlin values.
// 2. Sort them.
// 3. To transform a new value, find its "rank" in that sorted list.
// 4. Rank/10000 will give you a nearly perfect Uniform Distribution

func (s Setup) TestDistribution(
	testName string,
	noise noise.Noise,
) {
	s.T.Helper()
	expectedDeviation := internal.UniformDistribution
	noiseDeviation := internal.SampleNoiseDistribution(noise, 2000).StandardDeviation()
	const tolerance = 0.2
	if math.Abs(expectedDeviation-noiseDeviation) > tolerance {
		s.T.Errorf(
			"expected \"%s\" distribution to have deviation %.2f ± %.3f but got %v",
			testName,
			expectedDeviation,
			tolerance,
			noiseDeviation,
		)
	}
}
