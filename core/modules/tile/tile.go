package tile

type PosComponent struct {
	X, Y  int32
	Layer Layer
}

func NewPos(x, y int32, layer Layer) PosComponent {
	return PosComponent{x, y, layer}
}
