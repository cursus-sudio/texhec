package tileservice

import (
	"core/modules/tile"
	"engine/modules/relation"
	"engine/services/ecs"
)

type tool struct {
	tilePos relation.Service[tile.PosComponent]

	posArray ecs.ComponentsArray[tile.PosComponent]
}

func (t *tool) PosKey() relation.Service[tile.PosComponent] { return t.tilePos }

func (t *tool) Pos() ecs.ComponentsArray[tile.PosComponent] { return t.posArray }
