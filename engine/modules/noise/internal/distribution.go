package internal

import (
	"engine/modules/noise"
	"math"
	"sync"

	"github.com/go-gl/mathgl/mgl64"
)

var (
	UniformDistribution = 1. / math.Sqrt(12)
)

type Distribution []float64

func SampleNoiseDistribution(
	noise noise.Noise,
	samplesSqrt int,
) Distribution {
	mul := math.Pi * 1
	distribution := Distribution(make([]float64, samplesSqrt*samplesSqrt))
	wg := &sync.WaitGroup{}
	for x := range samplesSqrt {
		wg.Add(1)
		xCoord := float64(x) * mul
		go func() {
			defer wg.Done()
			for y := range samplesSqrt {
				yCoord := float64(y) * mul
				v := noise.Read(mgl64.Vec2{xCoord, yCoord})
				distribution[x*samplesSqrt+y] = v
			}
		}()
	}
	wg.Wait()

	return distribution
}

func (d Distribution) StandardDeviation() float64 {
	if len(d) == 0 {
		return 0
	}

	var sum float64
	for _, v := range d {
		sum += v
	}
	mean := sum / float64(len(d))

	var squaredDiffSum float64
	for _, v := range d {
		diff := v - mean
		squaredDiffSum += diff * diff
	}

	variance := squaredDiffSum / float64(len(d))
	return math.Sqrt(variance)
}
