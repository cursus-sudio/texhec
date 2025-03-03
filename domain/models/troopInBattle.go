package models

import (
	"domain/blueprints"
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

type TroopInBattle struct {
	models.ModelBase
	Blueprint   *blueprints.Troop
	TroopBudget vo.TroopBudget
	Amount      uint
	Dealt       vo.HP
}

func NewTroopInBattle(c ioc.Dic, troop *TroopInGame) TroopInBattle {
	return TroopInBattle{
		ModelBase:   models.NewBase(c),
		Blueprint:   troop.Blueprint,
		TroopBudget: troop.Blueprint.TroopBudget,
		Amount:      troop.Amount,
		Dealt:       vo.NoHP(),
	}
}

func (troop *TroopInBattle) StartTurn() {
	troop.TroopBudget = troop.Blueprint.TroopBudget
}

func (troop *TroopInBattle) Deal(dmg vo.HP) {
	var hp vo.HP
	for dmg.IsPositive() && troop.Amount > 0 {
		hp, dmg = troop.Blueprint.HP.Deal(dmg)
		if !hp.IsPositive() {
			troop.Amount -= 1
			continue
		}
		hp, dmg = hp.Deal(troop.Dealt)
		if !hp.IsPositive() {
			troop.Amount -= 1
			troop.Dealt = dmg
			continue
		}
		break
	}
}
