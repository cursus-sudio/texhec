package transform

type Size struct {
	X, Y, Z float32
}

func NewSize(x, y, z float32) Size {
	return Size{X: x, Y: y, Z: z}
}
