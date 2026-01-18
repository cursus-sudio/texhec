package tile

import "engine/modules/grid"

type Type uint8

func NewGrid(w, h grid.Coord) grid.SquareGridComponent[Type] {
	return grid.NewSquareGrid[Type](w, h)
}
