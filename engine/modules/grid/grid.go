package grid

import (
	"golang.org/x/exp/constraints"
)

type TileConstraint constraints.Unsigned

type Coord int16
type Index int

type SquareGridComponent[Tile TileConstraint] struct {
	width, height Coord
	grid          []Tile
}

func NewSquareGrid[Tile TileConstraint](w, h Coord) SquareGridComponent[Tile] {
	return SquareGridComponent[Tile]{
		width:  w,
		height: h,
		grid:   make([]Tile, int(w)*int(h)),
	}
}

// getters for consts
func (g *SquareGridComponent[Tile]) Width() Coord  { return g.width }
func (g *SquareGridComponent[Tile]) Height() Coord { return g.height }

// index and coord getters
func (g *SquareGridComponent[Tile]) GetIndex(x, y Coord) (Index, bool) {
	if x < 0 || y < 0 || x >= g.width || y >= g.height {
		return 0, false
	}
	return Index(x) + Index(y)*Index(g.width), true
}
func (g *SquareGridComponent[Tile]) GetCoords(index Index) (x, y Coord) {
	return Coord(index) % g.width, Coord(index) / g.width
}

// getters and setters tiles
func (g *SquareGridComponent[Tile]) GetTile(index Index) Tile       { return g.grid[index] }
func (g *SquareGridComponent[Tile]) SetTile(index Index, tile Tile) { g.grid[index] = tile }
