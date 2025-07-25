package transform

type Pos struct {
	X, Y, Z float32
}

func NewPos(x, y, z float32) Pos {
	return Pos{X: x, Y: y, Z: z}
}
