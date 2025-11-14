package tile

import "shared/services/ecs"

// per tile we want:
// - 1-3 textures to draw
// - 1 button with 1-3 different sets of buttons rendered at once

type Pos struct {
	X, Y int32
}

type Tile struct {
	DefID ecs.EntityID
}

type Building struct {
	DefID ecs.EntityID
}

type Unit struct {
	DefID ecs.EntityID
}

//

type HP uint32

// if not added then definition max hp is used
type Health struct {
	HP HP
}

type MaxHealth struct {
	MaxHP HP
}

// example:
// 1. soldier def has:
// - max health
// 2. soldier has:
// - unit
// - pos
// - health (optional)
