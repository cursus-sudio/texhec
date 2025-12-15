package tile

import (
	"engine/modules/relation"
)

type Tile interface {
	Tile() Interface
}

type Interface interface {
	TilePos() relation.EntityToKeyTool[PosComponent]
	ColliderPos() relation.EntityToKeyTool[ColliderPos]
}
