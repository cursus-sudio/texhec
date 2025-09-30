package tile

type TileTypeID uint16

type TilePos struct{ X, Y int }

type TileComponent struct {
	Pos  TilePos
	Type TileTypeID
}
