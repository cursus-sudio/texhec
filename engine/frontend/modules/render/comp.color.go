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

func DefaultColor() ColorComponent {
	return NewColor(mgl32.Vec4{1, 1, 1, 1})
}
