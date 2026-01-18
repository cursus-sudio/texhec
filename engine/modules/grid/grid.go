package grid

type TileConstraint interface {
	comparable
}

type Coord int16
type Coords struct{ X, Y Coord }

func (c *Coords) Coords() (X, Y Coord) {
	return c.X, c.Y
}

//

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
func (g *SquareGridComponent[Tile]) GetCoords(index Index) Coords {
	return Coords{
		X: Coord(index) % g.width,
		Y: Coord(index) / g.width,
	}
}
func (g *SquareGridComponent[Tile]) GetTiles() []Tile {
	tiles := make([]Tile, len(g.grid))
	copy(tiles, g.grid)
	return tiles
}
func (g *SquareGridComponent[Tile]) GetLastIndex() Index {
	return Index(g.width) * Index(g.height)
}

// getters and setters tiles
func (g *SquareGridComponent[Tile]) GetTile(index Index) Tile       { return g.grid[index] }
func (g *SquareGridComponent[Tile]) SetTile(index Index, tile Tile) { g.grid[index] = tile }
