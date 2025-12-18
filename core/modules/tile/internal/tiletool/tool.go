package tiletool

import (
	"core/modules/tile"
	"engine/modules/relation"
	"engine/services/ecs"
)

type tool struct {
	tilePos     relation.EntityToKeyTool[tile.PosComponent]
	colliderPos relation.EntityToKeyTool[tile.ColliderPos]

	posArray ecs.ComponentsArray[tile.PosComponent]
}

func (t *tool) Tile() tile.Interface { return t }

func (t *tool) PosKey() relation.EntityToKeyTool[tile.PosComponent]        { return t.tilePos }
func (t *tool) ColliderPosKey() relation.EntityToKeyTool[tile.ColliderPos] { return t.colliderPos }

func (t *tool) Pos() ecs.ComponentsArray[tile.PosComponent] { return t.posArray }
