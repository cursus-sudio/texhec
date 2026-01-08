package tile

import (
	"golang.org/x/exp/constraints"
)

type PosComponent struct {
	X, Y  float32
	Layer Layer
}

func NewPos[Number constraints.Float | constraints.Integer](x, y Number, layer Layer) PosComponent {
	return PosComponent{float32(x), float32(y), layer}
}
