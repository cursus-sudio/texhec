package tile

import "engine/services/ecs"

type System ecs.SystemRegister
type SystemRenderer ecs.SystemRegister

type TileClickEvent struct{ Tile ecs.EntityID }

func NewTileClickEvent(t ecs.EntityID) TileClickEvent { return TileClickEvent{t} }
