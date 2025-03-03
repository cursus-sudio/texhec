package models

import (
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

type MapInBattle struct {
	models.ModelBase
	Size  vo.SizeInBattle
	Tiles map[vo.PositionInBattle]*TileInBattle
}

// SERVICE
type MapInBattleGenerator interface {
	// TODO
	New(size vo.SizeInBattle, squads ...*TroopArmyInGame) *MapInBattle
}

func NewMapInBattle(c ioc.Dic, size vo.SizeInBattle, armies ...*TroopArmyInGame) *MapInBattle {
	return ioc.Get[MapInBattleGenerator](c).New(size, armies...)
}
