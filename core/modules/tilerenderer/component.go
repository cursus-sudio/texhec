package tilerenderer

type TilePosComponent struct{ X, Y, Z int32 }

func NewTilePos(x, y, z int32) TilePosComponent { return TilePosComponent{x, y, z} }

type TileTypeComponent struct {
	Type uint32
}
