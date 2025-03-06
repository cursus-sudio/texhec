package battle

import (
	"domain/common/models"
	"domain/states/game/gamedefinitions"
	"domain/states/kernel/vo"
)

type Troop struct {
	models.ModelBase
	Definition  gamedefinitions.Troop
	TroopBudget vo.TroopBudget
	Amount      uint
	Dealt       vo.Hp
}
