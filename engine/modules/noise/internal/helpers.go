package internal

import (
	"engine/modules/noise"
	"math"
	"sync"
	"sync/atomic"

	"github.com/go-gl/mathgl/mgl64"
)

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func dot(v1, v2 mgl64.Vec2) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}

func fract(x float64) float64 {
	return x - math.Floor(x)
}

var c00 = mgl64.Vec2{0, 0}
var c10 = mgl64.Vec2{1, 0}
var c01 = mgl64.Vec2{0, 1}
var c11 = mgl64.Vec2{1, 1}

//

func MapRange(v, inMin, inMax, outMin, outMax float64) float64 {
	return outMin + (v-inMin)*(outMax-outMin)/(inMax-inMin)
}

// to get values create example 3 segment chart and pass here 3 and procentage value of main segment
func standardDeviation(segmentsCount, valueInMainSegment float64) float64 {
	return .5 / (segmentsCount * math.Sqrt(2) * math.Erfinv(valueInMainSegment))
}

// cumulative distribution function
// s is StandardDeviation
func cdf(x, s float64) float64 {
	return 0.5 * (1 + math.Erf((x-0.5)/(s*math.Sqrt(2))))
}

func CalculateDistribution(
	noise noise.Noise,
	samplesSqrt int,
) [3]float64 {
	mul := math.Pi
	count := [3]int64{}
	wg := &sync.WaitGroup{}
	for x := range samplesSqrt {
		wg.Add(1)
		xCoord := float64(x) * mul
		go func() {
			defer wg.Done()
			for y := range samplesSqrt {
				yCoord := float64(y) * mul
				v := noise.Read(mgl64.Vec2{xCoord, yCoord})
				atomic.AddInt64(&count[min(int(v*3), 2)], 1)
			}
		}()
	}
	wg.Wait()
	res := [3]float64{}
	for i := range 3 {
		res[i] = float64(count[i]) / float64(samplesSqrt*samplesSqrt)
	}

	return res
}
