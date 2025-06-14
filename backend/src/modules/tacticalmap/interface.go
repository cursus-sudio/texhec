package tacticalmap

type Pos struct {
	X, Y int
}

type Tile struct {
	Pos Pos
}

type CreateArgs struct {
	Tiles []Tile
}

type DestroyArgs struct {
	Tiles []Tile
}

type CreateListener func(tiles []Tile)
type DestroyListener func(tiles []Tile)

type TacticalMap interface {
	Create(CreateArgs) error
	Destroy(DestroyArgs) error
	GetMap() ([]Tile, error)
	OnCreate(CreateListener)
	OnDestroy(DestroyListener)
}
