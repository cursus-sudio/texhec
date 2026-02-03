package internal

import (
	"engine/modules/noise"
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

func valueHash(seed uint64, p mgl64.Vec2) float64 {
	s := float64(seed) / math.MaxUint64
	p = mgl64.Vec2{
		fract(p.X()*123.34 + s*10.123),
		fract(p.Y()*456.21 + s*20.456),
	}
	shift := s * 100.0
	dotVal := p.Dot(p.Add(mgl64.Vec2{45.32 + shift, 45.32 + shift}))
	p = p.Add(mgl64.Vec2{dotVal, dotVal})

	return fract(p.X() * p.Y())
}

func valueInterpolate(t float64) float64 {
	return t * t * (3 - 2*t)
}

func NewValueNoise(seed uint64, layer noise.LayerConfig) Noise {
	return NewNoise(func(coords mgl64.Vec2) float64 {
		coords = coords.Add(layer.Offset)
		coords = coords.Mul(layer.CellSize)

		i := mgl64.Vec2{math.Floor(coords.X()), math.Floor(coords.Y())}
		f := mgl64.Vec2{fract(coords.X()), fract(coords.Y())}

		a := valueHash(seed, i.Add(c00))
		b := valueHash(seed, i.Add(c10))
		c := valueHash(seed, i.Add(c01))
		d := valueHash(seed, i.Add(c11))

		ux := valueInterpolate(f.X())
		uy := valueInterpolate(f.Y())

		res := lerp(
			lerp(a, b, ux),
			lerp(c, d, ux),
			uy,
		)

		return res * layer.ValueMultiplier
	})
}

// value noise
// float hash(vec2 p) {
//     p = fract(p * vec2(123.34, 456.21));
//     p += dot(p, p + 45.32);
//     return fract(p.x * p.y);
// }
// vec2 interpolate(vec2 p) {
//     return p * p * (3 - 2 * p);
// }
// float value_noise(vec2 uv) { // normalized <0, 1>
//     vec2 i = floor(uv);
//     vec2 f = fract(uv);
//
//     float a = hash(i);
//     float b = hash(i + vec2(1, 0));
//     float c = hash(i + vec2(0, 1));
//     float d = hash(i + vec2(1, 1));
//
//     vec2 u = interpolate(f);
//
//     return mix(a, b, u.x) +
//         (c - a) * u.y * (1 - u.x) +
//         (d - b) * u.x * u.y;
// }
