package models

import (
	"domain/blueprints"
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/null"
)

// SERVICE
type TroopInGameErrors struct {
	TroopIsAlreadyFull,
	CannotJoinDifferentTroop error
}

type TroopInGame struct {
	models.ModelBase
	Blueprint *blueprints.Troop
	Budget    vo.TroopBudget
	Amount    uint
}

func NewTroopInGame(c ioc.Dic, blueprint *blueprints.Troop, currentBudget vo.TroopBudget, amount uint) *TroopInGame {
	return &TroopInGame{
		ModelBase: models.NewBase(c),
		Blueprint: blueprint,
		Budget:    currentBudget,
		Amount:    amount,
	}
}

func (troop *TroopInGame) StartTurn(c ioc.Dic) {
	troop.Budget = troop.Blueprint.TroopBudget
}

func (troop *TroopInGame) CanJoin(c ioc.Dic, joined *TroopInGame) (after *TroopInGame, overlow null.Nullable[*TroopInGame], err error) {
	errors := ioc.Get[TroopInGameErrors](c)
	if troop.Blueprint != joined.Blueprint {
		return nil, null.Null[*TroopInGame](), errors.CannotJoinDifferentTroop
	}

	if troop.Amount == troop.Blueprint.MaxAmount {
		return nil, null.Null[*TroopInGame](), errors.TroopIsAlreadyFull
	}

	sum := troop.Amount + joined.Amount
	minBudget := vo.MinTroopBudget(troop.Budget, joined.Budget)
	if sum <= troop.Blueprint.MaxAmount {
		return NewTroopInGame(c, troop.Blueprint, minBudget, sum), null.Null[*TroopInGame](), nil
	}

	return NewTroopInGame(c, troop.Blueprint, minBudget, troop.Blueprint.MaxAmount),
		null.New(NewTroopInGame(c, troop.Blueprint, joined.Budget, troop.Blueprint.MaxAmount-sum)),
		nil
}

func (troop *TroopInGame) Join(c ioc.Dic, joined *TroopInGame) (overflow null.Nullable[*TroopInGame], err error) {
	changes, o, err := troop.CanJoin(c, joined)
	if err != nil {
		return o, err
	}
	troop.Amount = changes.Amount
	troop.Budget = changes.Budget
	return o, err
}
