package gameobjects

import (
	"domain/common/models"
	"domain/states/game/gamedefinitions"
	"domain/states/game/gamevo"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/null"
)

type Tile struct {
	models.ModelBase
	pos           gamevo.Pos
	definition    gamedefinitions.Tile
	resource      null.Nullable[gamedefinitions.Resource]
	building      null.Nullable[*Building]
	troop         null.Nullable[*TroopArmy]
	troopAttacker null.Nullable[*TroopArmy]
}

func NewTile(
	c ioc.Dic,
	base models.ModelBase,
	pos gamevo.Pos,
	definition gamedefinitions.Tile,
	resource null.Nullable[gamedefinitions.Resource],
	building null.Nullable[*Building],
	troop null.Nullable[*TroopArmy],
	troopAttacker null.Nullable[*TroopArmy],
) Tile {
	return Tile{
		ModelBase:     base,
		pos:           pos,
		definition:    definition,
		resource:      resource,
		building:      building,
		troop:         troop,
		troopAttacker: troopAttacker,
	}
}

// func
