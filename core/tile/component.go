package tile

type TilePos struct{ X, Y int32 }

type TileComponent struct {
	Pos  TilePos
	Type uint32
}
