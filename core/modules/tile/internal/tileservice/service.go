package tileservice

import (
	"core/modules/tile"
	"engine/modules/relation"
	"engine/services/ecs"
)

type service struct {
	tilePos relation.Service[tile.PosComponent]

	posArray ecs.ComponentsArray[tile.PosComponent]
}

func (t *service) PosKey() relation.Service[tile.PosComponent] { return t.tilePos }

func (t *service) Pos() ecs.ComponentsArray[tile.PosComponent] { return t.posArray }
