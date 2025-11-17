package tile

const (
	GroundLayer uint8 = iota
	BuildingLayer
	UnitLayer
)

type PosComponent struct {
	X, Y  int32
	Layer uint8
}

func NewPos(x, y int32, layer uint8) PosComponent {
	return PosComponent{x, y, layer}
}
