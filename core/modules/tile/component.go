package tile

type TilePos struct{ X, Y, Z int32 }

func NewTilePos(x, y, z int32) TilePos { return TilePos{x, y, z} }

type TileComponent struct {
	Pos  TilePos
	Type uint32
}
