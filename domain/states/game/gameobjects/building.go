package gameobjects

import (
	"domain/common/models"
	"domain/states/game/gamedefinitions"
	"domain/states/kernel/vo"

	"github.com/ogiusek/ioc"
)

// SERVICE
type BuildingErrors interface {
	BuildingHpCannotExceedMaxHp(building Building, hp vo.Hp) error
}

type Building struct {
	models.ModelBase
	definition gamedefinitions.Building
	hp         vo.Hp
}

func NewBuilding(c ioc.Dic, base models.ModelBase, definition gamedefinitions.Building) *Building {
	return &Building{
		ModelBase:  base,
		definition: definition,
		// TODO
		// hp:         definition.MaxHp(),
	}
}

func (building *Building) Definition() gamedefinitions.Building {
	return building.definition
}

func (building *Building) Hp() vo.Hp {
	return building.hp
}

func (building *Building) StartTurn(c ioc.Dic) {
	// TODO
	// trigger listeners
}

func (building *Building) SetHp(c ioc.Dic, hp vo.Hp) error {
	if hp.IsGreatherThan(building.definition.MaxHp()) {
		errors := ioc.Get[BuildingErrors](c)
		return errors.BuildingHpCannotExceedMaxHp(*building, hp)
	}
	building.hp = hp
	return nil
}
