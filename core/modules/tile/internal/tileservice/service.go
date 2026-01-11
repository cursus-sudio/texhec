package tileservice

import (
	"core/modules/tile"
	"engine/modules/relation"
	"engine/services/ecs"
)

type service struct {
	TilePos  relation.Service[tile.PosComponent]
	PosArray ecs.ComponentsArray[tile.PosComponent]
}

func (t *service) PosKey() relation.Service[tile.PosComponent] { return t.TilePos }

func (t *service) Pos() ecs.ComponentsArray[tile.PosComponent] { return t.PosArray }
