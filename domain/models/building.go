package models

import (
	"domain/blueprints"
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

type Building struct {
	models.ModelBase
	Blueprint *blueprints.Building
	HP        vo.HP
	Active    bool
}

func NewBuilding(c ioc.Dic, blueprint *blueprints.Building) *Building {
	return &Building{
		ModelBase: models.NewBase(c),
		Blueprint: blueprint,
		HP:        blueprint.HP,
		Active:    true,
	}
}
