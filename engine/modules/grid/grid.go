package grid

import "golang.org/x/exp/constraints"

type Coord int16
type Index int

type SquareGridComponent[TileType constraints.Unsigned] struct {
	width, height Coord
	grid          []TileType
}

func NewSquareGrid[TileType constraints.Unsigned](w, h Coord) SquareGridComponent[TileType] {
	return SquareGridComponent[TileType]{
		width:  w,
		height: h,
		grid:   make([]TileType, int(w)*int(h)),
	}
}

// getters for consts
func (g *SquareGridComponent[TileType]) Width() Coord  { return g.width }
func (g *SquareGridComponent[TileType]) Height() Coord { return g.height }

// index and coord getters
func (g *SquareGridComponent[TileType]) GetIndex(x, y Coord) (Index, bool) {
	if x < 0 || y < 0 || x >= g.width || y >= g.height {
		return 0, false
	}
	return Index(x) + Index(y)*Index(g.width), true
}
func (g *SquareGridComponent[TileType]) GetCoords(index Index) (x, y Coord) {
	return Coord(index) % g.width, Coord(index) / g.width
}

// getters and setters tiles
func (g *SquareGridComponent[TileType]) GetTile(index Index) TileType       { return g.grid[index] }
func (g *SquareGridComponent[TileType]) SetTile(index Index, tile TileType) { g.grid[index] = tile }
