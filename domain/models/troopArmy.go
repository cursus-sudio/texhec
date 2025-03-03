package models

import (
	"domain/common/models"

	"github.com/ogiusek/ioc"
)

// SERVICE
type TroopArmyLimits struct {
	MaxSlots int
}

type TroopArmyInGame struct {
	models.ModelBase
	PlayerId models.ModelId
	Troops   []*TroopInGame
}

func NewTroopArmyInGame(c ioc.Dic, playerId models.ModelId, troops []*TroopInGame) (TroopArmyInGame, error) {
	limits := ioc.Get[TroopArmyLimits](c)
	if limits.MaxSlots < len(troops) {
		return TroopArmyInGame{}, nil
	}

	return TroopArmyInGame{
		ModelBase: models.NewBase(c),
		PlayerId:  playerId,
		Troops:    troops,
	}, nil
}

func (army *TroopArmyInGame) StartTurn(c ioc.Dic) {
	for _, troop := range army.Troops {
		troop.StartTurn(c)
	}
}

func (army *TroopArmyInGame) CanJoin(c ioc.Dic, joined *TroopArmyInGame) error {
	// TODO
	return nil
}

func (army *TroopArmyInGame) Join(c ioc.Dic, joined *TroopArmyInGame) error {
	// TODO
	return nil
}
