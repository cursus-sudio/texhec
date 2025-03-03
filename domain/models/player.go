package models

import (
	"domain/blueprints"
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

type Player struct {
	models.ModelBase
	User            *User
	DiscoveredTiles map[vo.PositionInGame]*DiscoveredTile
	Fraction        *blueprints.Fraction
	BrainBuilding   *Building
	Bugdet          vo.Budget
	Buildings       []*Building
	Troops          []*TroopArmyInGame
}

func NewPlayer(
	c ioc.Dic,
	user *User,
	fraction *blueprints.Fraction,
	brainBuilding *Building,
	defaultInitialBudget vo.Budget,
) *Player {
	return &Player{
		ModelBase:       models.NewBase(c),
		User:            user,
		DiscoveredTiles: map[vo.PositionInGame]*DiscoveredTile{},
		Fraction:        fraction,
		BrainBuilding:   brainBuilding,
		Bugdet:          defaultInitialBudget.Multiply(fraction.InitialBudgetMultiplier),
		Buildings:       []*Building{},
		Troops:          []*TroopArmyInGame{},
	}
}
