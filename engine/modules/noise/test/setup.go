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
	setup.Layer = noise.NewLayer(10, 1)
	return setup
}

func (s Setup) TestDistribution(testName string, noise noise.Noise) {
	s.T.Helper()
	distribution := internal.CalculateDistribution(noise, 1000)
	for _, value := range distribution {
		const target = 1. / 3
		const tolerance = .02
		if math.Abs(target-value) > tolerance {
			s.T.Errorf(
				"expected \"%s\" distribution to have only values %.2f ± %.2f but got %v",
				testName,
				target,
				tolerance,
				distribution,
			)
			return
		}
	}
}
