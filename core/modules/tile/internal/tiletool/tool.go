package tiletool

import (
	"core/modules/tile"
	"engine/modules/relation"
)

type tool struct {
	tilePos     relation.EntityToKeyTool[tile.PosComponent]
	colliderPos relation.EntityToKeyTool[tile.ColliderPos]
}

func (t *tool) TilePos() relation.EntityToKeyTool[tile.PosComponent]    { return t.tilePos }
func (t *tool) ColliderPos() relation.EntityToKeyTool[tile.ColliderPos] { return t.colliderPos }
