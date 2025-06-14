package tacticalmap

type Repo interface {
	Create() error
	Delete() error
	GetMap() ([]Tile, error)
	OnCreate(CreateListener)
	OnDelete(DestroyListener)
}
