package render

import "github.com/go-gl/mathgl/mgl32"

// normalized color applied
// default is (1, 1, 1, 1)
type ColorComponent struct {
	Color mgl32.Vec4
}

func NewColor(color mgl32.Vec4) ColorComponent {
	return ColorComponent{
		Color: color,
	}
}

func (c1 ColorComponent) Lerp(c2 ColorComponent, mix32 float32) ColorComponent {
	return ColorComponent{c1.Color.Mul(1 - mix32).Add(c2.Color.Mul(mix32))}
}
