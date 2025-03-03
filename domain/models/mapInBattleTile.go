package models

import (
	"domain/blueprints"
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/null"
)

type TileInBattle struct {
	models.ModelBase
	Position  vo.PositionInBattle
	Blueprint *blueprints.TileInBattle
	Troop     null.Nullable[*TroopInBattle]
}

func NewTileInBattle(c ioc.Dic, position vo.PositionInBattle, blueprint *blueprints.TileInBattle) TileInBattle {
	return TileInBattle{
		ModelBase: models.NewBase(c),
		Position:  position,
		Blueprint: blueprint,
		Troop:     null.Null[*TroopInBattle](),
	}
}
