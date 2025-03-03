package blueprints

import (
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

type Troop struct {
	models.ModelBase
	models.ModelDescription
	Upgrades       []*Troop
	TroopBudget    vo.TroopBudget
	MaxAmount      uint
	HP             vo.HP
	SizeInBattle   vo.SizeInBattle
	Cost           vo.Budget
	SkillsInBattle []*SkillInBattle
	SkillsInGame   []*SkillInGame
	Builds         []*Building
	BuildCapacity  vo.Budget
}

func NewTroop(
	c ioc.Dic,
	desc models.ModelDescription,
	upgrades []*Troop,
	troopBudget vo.TroopBudget,
	maxAmount uint,
	hp vo.HP,
	sizeInBattle vo.SizeInBattle,
	cost vo.Budget,
	skillsInBattle []*SkillInBattle,
	skillsInGame []*SkillInGame,
	builds []*Building,
	buildCapacity vo.Budget,
) Troop {
	return Troop{
		ModelBase:        models.NewBase(c),
		ModelDescription: desc,
		Upgrades:         upgrades,
		TroopBudget:      troopBudget,
		MaxAmount:        maxAmount,
		HP:               hp,
		SizeInBattle:     sizeInBattle,
		Cost:             cost,
		SkillsInBattle:   skillsInBattle,
		SkillsInGame:     skillsInGame,
		Builds:           builds,
		BuildCapacity:    buildCapacity,
	}
}
