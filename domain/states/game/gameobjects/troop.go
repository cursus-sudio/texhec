package gameobjects

import (
	"domain/common/gameobject"
	"domain/common/models"
	"domain/states/game/gamedefinitions"
	"domain/states/kernel/vo"

	"github.com/ogiusek/ioc"
)

// SERVICE
type TroopErrors interface {
	CannotExceedMaxAmount(troop Troop, amount uint) error
}

type Troop struct {
	models.ModelBase
	metadata   gameobject.Metadata
	definition gamedefinitions.Troop
	budget     vo.TroopBudget
	amount     uint
}

func NewTroop(c ioc.Dic, base models.ModelBase, definition gamedefinitions.Troop, amount uint) Troop {
	troop := Troop{
		ModelBase:  base,
		metadata:   definition.Metadata(),
		definition: definition,
		budget:     definition.TroopBudget(),
		amount:     amount,
	}
	return troop
}

func (troop *Troop) Definition() gamedefinitions.Troop {
	return troop.definition
}

func (troop *Troop) Budget() vo.TroopBudget {
	return troop.budget
}

func (troop *Troop) Amount() uint {
	return troop.amount
}

// func (troop *Troop) Metadata()

func (troop *Troop) StartTurn() {
	troop.budget = troop.definition.TroopBudget()
}

func (troop *Troop) Pay(c ioc.Dic, pay *vo.TroopBudget) error {
	res, err := troop.budget.Pay(c, pay)
	if err != nil {
		return err
	}
	troop.budget = res
	return nil
}

func (troop *Troop) SetAmount(c ioc.Dic, amount uint) error {
	if troop.definition.MaxAmount() < amount {
		errors := ioc.Get[TroopErrors](c)
		return errors.CannotExceedMaxAmount(*troop, amount)
	}

	return nil
}

// func (troop *Troop)
