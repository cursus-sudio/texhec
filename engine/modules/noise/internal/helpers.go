package internal

import (
	"math"

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
// Standard Deviation
func standardDeviation(segmentsCount, valueInMainSegment float64) float64 {
	return .5 / (segmentsCount * math.Sqrt(2) * math.Erfinv(valueInMainSegment))
	// return .5 / (segmentsCount * math.Sqrt(2) * invErf(valueInMainSegment))
}

// s is StandardDeviation
func cdf(x, s float64) float64 {
	return 0.5 * (1 + math.Erf((x-0.5)/(s*math.Sqrt(2))))
}

// func invErf(y float64) float64 {
// 	if y < -1 || y > 1 {
// 		return math.NaN()
// 	}
// 	if y == 0 {
// 		return 0
// 	}
//
// 	// Quick approximation for the range you need
// 	// For better precision, use a dedicated library like gonum
// 	a := 0.147
// 	term1 := 2/(math.Pi*a) + math.Log(1-y*y)/2
// 	term2 := math.Log(1-y*y) / a
//
// 	part1 := math.Sqrt(term1*term1 - term2)
// 	return math.Copysign(math.Sqrt(part1-term1), y)
// }
