package tile

import (
	"engine/modules/relation"
)

type Tool interface {
	TilePos() relation.EntityToKeyTool[PosComponent]
	ColliderPos() relation.EntityToKeyTool[ColliderPos]
}
