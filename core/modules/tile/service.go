package tile

import (
	"engine/modules/grid"
	"engine/modules/transform"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/ecs"
)

type ID uint8

func NewGrid(w, h grid.Coord) grid.SquareGridComponent[ID] {
	return grid.NewSquareGrid[ID](w, h)
}

//

type Service interface {
	Grid() ecs.ComponentsArray[grid.SquareGridComponent[ID]]
	GetPos(grid.Coords) transform.PosComponent
	GetTileSize() transform.SizeComponent
}

type TileAssets interface {
	AddType(addedAssets datastructures.SparseArray[ID, assets.AssetID])
}

//

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
