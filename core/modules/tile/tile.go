package tile

// per tile we want:
// - 1-3 textures to draw
// - 1 button with 1-3 different sets of buttons rendered at once

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

//

// example:
// 1. soldier def has:
// - max health
// 2. soldier has:
// - unit
// - pos
// - health (optional)

// definitions can have:
// - max hp
// - cost
// - recruits
// - round profit
// - attacks
// - defense (optional, armor, immunities etc)
// -

// attack:
// - min, max range
// - damage
