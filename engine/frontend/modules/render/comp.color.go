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

func (c1 ColorComponent) Blend(c2 ColorComponent, mix32 float32) ColorComponent {
	invMix32 := 1.0 - mix32
	blendedColor := mgl32.Vec4{
		c1.Color[0]*invMix32 + c2.Color[0]*mix32,
		c1.Color[1]*invMix32 + c2.Color[1]*mix32,
		c1.Color[2]*invMix32 + c2.Color[2]*mix32,
		c1.Color[3]*invMix32 + c2.Color[3]*mix32,
	}

	return ColorComponent{
		Color: blendedColor,
	}
}

func DefaultColor() ColorComponent {
	return NewColor(mgl32.Vec4{1, 1, 1, 1})
}
