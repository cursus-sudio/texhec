package tile

import (
	"engine/modules/grid"
	"engine/services/ecs"
)

type System ecs.SystemRegister
type SystemRenderer ecs.SystemRegister

type TileClickEvent struct {
	Grid ecs.EntityID
	Tile grid.Index
}

func NewTileClickEvent(
	grid ecs.EntityID,
	tile grid.Index,
) any {
	return TileClickEvent{grid, tile}
}
