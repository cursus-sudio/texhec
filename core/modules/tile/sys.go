package tile

import (
	"engine/services/ecs"
)

type System ecs.SystemRegister[World]
type SystemRenderer ecs.SystemRegister[World]

type TileClickEvent struct{ Tile PosComponent }

func NewTileClickEvent(t PosComponent) TileClickEvent { return TileClickEvent{t} }
