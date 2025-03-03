package game

import (
	"domain/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

// SERVICE
type BattleDefaults struct {
	MapSize vo.SizeInBattle
}

type Battle struct {
	Position vo.PositionInGame
	Armies   []*models.TroopArmyInGame
	Turn     int
	Map      *models.MapInBattle
}

func NewBattle(c ioc.Dic, position vo.PositionInGame, armies ...*models.TroopArmyInGame) *Battle {
	defaults := ioc.Get[BattleDefaults](c)

	return &Battle{
		Position: position,
		Armies:   armies,
		Turn:     0,
		Map:      models.NewMapInBattle(c, defaults.MapSize, armies...),
	}
}
