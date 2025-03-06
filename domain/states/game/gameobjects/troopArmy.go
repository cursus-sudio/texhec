package gameobjects

import (
	"domain/common/models"
	"errors"

	"github.com/ogiusek/ioc"
)

// SERVICE
type TroopArmyErrors interface {
	ExceededMaxSlots(maxSlots uint, current uint) error
}

// SERVICE
type TroopArmyConstants struct {
	troopSlots uint
}

func NewTroopArmyConstants(troopSlots uint) (TroopArmyConstants, error) {
	if troopSlots == 0 {
		// this error is not injected because this is configuration error
		return TroopArmyConstants{}, errors.New("troop army slots have to be positive")
	}
	return TroopArmyConstants{
		troopSlots: troopSlots,
	}, nil
}

//

type TroopArmy struct {
	models.ModelBase
	playerId models.ModelId
	troops   map[models.ModelId]*Troop
}

func NewTroopArmy(c ioc.Dic, base models.ModelBase, playerId models.ModelId, troops map[models.ModelId]*Troop) (TroopArmy, error) {
	errors := ioc.Get[TroopArmyErrors](c)
	constraints := ioc.Get[TroopArmyConstants](c)

	troopsLen := len(troops)
	if troopsLen > int(constraints.troopSlots) {
		return TroopArmy{}, errors.ExceededMaxSlots(constraints.troopSlots, uint(troopsLen))
	}

	return TroopArmy{
		ModelBase: base,
		playerId:  playerId,
		troops:    troops,
	}, nil
}

func (army *TroopArmy) Troops() map[models.ModelId]Troop {
	troops := make(map[models.ModelId]Troop, len(army.troops))
	for id, troop := range army.troops {
		troops[id] = *troop
	}
	return troops
}

func (army *TroopArmy) UpdateTroop(troop *Troop) error {
	// TODO
	return nil
}

func (army *TroopArmy) RemoveTroop(troopId models.ModelId) {
	// TODO
	// delete()
}

func (army *TroopArmy) AddTroop(troop *Troop) error {
	// TODO
	return nil
}
