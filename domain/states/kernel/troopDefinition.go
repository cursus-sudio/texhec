package kernel

import (
	"domain/common/models"
	"domain/states/kernel/vo"

	"github.com/ogiusek/ioc"
)

type TroopDefinition struct {
	models.ModelBase
	models.ModelDescription
	maxHp       vo.Hp
	cost        vo.Budget
	troopBudget vo.TroopBudget
	maxAmount   uint
}

func NewTroop(c ioc.Dic, base models.ModelBase, desc models.ModelDescription, maxHp vo.Hp, cost vo.Budget, troopBudget vo.TroopBudget, maxAmount uint) TroopDefinition {
	return TroopDefinition{
		ModelBase:        base,
		ModelDescription: desc,
		maxHp:            maxHp,
		cost:             cost,
		troopBudget:      troopBudget,
		maxAmount:        maxAmount,
	}
}

func (troop *TroopDefinition) MaxHp() vo.Hp {
	return troop.maxHp
}

func (troop *TroopDefinition) Cost() vo.Budget {
	return troop.cost
}

func (troop *TroopDefinition) TroopBudget() vo.TroopBudget {
	return troop.troopBudget
}

func (troop *TroopDefinition) MaxAmount() uint {
	return troop.maxAmount
}
