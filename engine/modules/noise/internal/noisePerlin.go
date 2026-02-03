package internal

import (
	"engine/modules/noise"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

func perlinHash(seed uint64, c mgl64.Vec2) mgl64.Vec2 {
	sLow := uint32(seed)
	sHigh := uint32(seed >> 32)

	x := uint32(int32(c.X())) ^ sLow
	y := uint32(int32(c.Y())) ^ sHigh

	const mult uint32 = 1664525
	const increment uint32 = 1013904223

	x = x*mult + increment
	y = y*mult + increment

	x += y * mult * sLow
	y += x * mult * sHigh

	x ^= (x >> 16)
	y ^= (y >> 16)

	x += y * mult
	y += x * mult

	x ^= (x >> 16)
	y ^= (y >> 16)

	const maxUint float64 = 4294967295.0

	resX := -1.0 + 2.0*(float64(x)/maxUint)
	resY := -1.0 + 2.0*(float64(y)/maxUint)

	return mgl64.Vec2{resX, resY}
}

func perlinInterpolate(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func NewPerlinNoise(seed uint64, layer noise.LayerConfig) Noise {
	return NewNoise(func(coords mgl64.Vec2) float64 {
		coords = coords.Add(layer.Offset)
		coords = coords.Mul(layer.CellSize)
		i := mgl64.Vec2{math.Floor(coords.X()), math.Floor(coords.Y())}
		f := coords.Sub(i)

		h00 := perlinHash(seed, i.Add(c00))
		a := dot(h00, f.Sub(c00))

		h10 := perlinHash(seed, i.Add(c10))
		b := dot(h10, f.Sub(c10))

		h01 := perlinHash(seed, i.Add(c01))
		c := dot(h01, f.Sub(c01))

		h11 := perlinHash(seed, i.Add(c11))
		d := dot(h11, f.Sub(c11))

		ux := perlinInterpolate(f.X())
		uy := perlinInterpolate(f.Y())

		res := lerp(
			lerp(a, b, ux),
			lerp(c, d, ux),
			uy,
		)

		res = res*0.5 + 0.5

		// stretch (.1, .9) to (0, 1) and clamp the rest
		res = (res - .1) / (.9 - .1)
		res = mgl64.Clamp(res, 0, 1)

		return res * layer.ValueMultiplier
	})
}

// glsl original
// vec2 hash(vec2 p) {
//     p = vec2(dot(p, vec2(127.1, 311.7)), dot(p, vec2(269.5, 183.3)));
//     return -1.0 + 2.0 * fract(sin(p) * 43758.5453123);
// }
//
// // Quintic interpolation curve (smoother than Hermite/Smoothstep)
// // Formula: 6t^5 - 15t^4 + 10t^3
// vec2 interpolate(vec2 p) {
//     return p * p * p * (p * (p * 6 - 15) + 10);
// }
//
// float perlin_noise(vec2 p) { // normalzied <0, 1>
//     vec2 i = floor(p);
//     vec2 f = fract(p);
//
//     float a = dot(hash(i + vec2(0, 0)), f - vec2(0, 0));
//     float b = dot(hash(i + vec2(1, 0)), f - vec2(1, 0));
//     float c = dot(hash(i + vec2(0, 1)), f - vec2(0, 1));
//     float d = dot(hash(i + vec2(1, 1)), f - vec2(1, 1));
//
//     vec2 u = interpolate(f);
//
//     float result = mix(mix(a, b, u.x), mix(c, d, u.x), u.y);
//     return result * .5 + .5;
// }

//
//
//

// old hash
// pX := c.X()*127.1 + c.Y()*311.7
// pY := c.X()*269.5 + c.Y()*183.3
//
// // sin(p) * 43758.5453123
// resX := math.Sin(pX) * 43758.5453123
// resY := math.Sin(pY) * 43758.5453123
//
// // -1.0 + 2.0 * fract(...)
// _, fracX := math.Modf(resX)
// _, fracY := math.Modf(resY)
//
// // Handling negative Modf results to ensure they stay in [0, 1]
// if fracX < 0 {
// 	fracX += 1.0
// }
// if fracY < 0 {
// 	fracY += 1.0
// }
//
// return mgl64.Vec2{
// 	-1.0 + 2.0*fracX,
// 	-1.0 + 2.0*fracY,
// }
