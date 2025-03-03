package blueprints

import (
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

type Fraction struct {
	models.ModelBase
	models.ModelDescription
	SpawnBiom                *Biom
	BrainBlueprint           *Building
	InitialResearch          *Research
	InitialBudgetMultiplier  vo.BudgetMultiplier
	BuildingIncomeMultiplier vo.BudgetMultiplier
	BuildingCostMultiplier   vo.BudgetMultiplier
	TroopCostMultiplier      vo.BudgetMultiplier
	TroopBudgetMultiplier    vo.TroopBudgetMultiplier
}

func NewFraction(
	c ioc.Dic,
	desc models.ModelDescription,
	spawnBiom *Biom,
	brainBuilding *Building,
	initialResearch *Research,
	initialBudgetMultiplier vo.BudgetMultiplier,
	buildingIncomeMultiplier vo.BudgetMultiplier,
	buildingCostMultiplier vo.BudgetMultiplier,
	troopCostMultiplier vo.BudgetMultiplier,
	troopBudgetMultiplier vo.TroopBudgetMultiplier,
) Fraction {
	return Fraction{
		ModelBase:                models.NewBase(c),
		ModelDescription:         desc,
		SpawnBiom:                spawnBiom,
		BrainBlueprint:           brainBuilding,
		InitialResearch:          initialResearch,
		InitialBudgetMultiplier:  initialBudgetMultiplier,
		BuildingIncomeMultiplier: buildingIncomeMultiplier,
		BuildingCostMultiplier:   buildingCostMultiplier,
		TroopCostMultiplier:      troopCostMultiplier,
		TroopBudgetMultiplier:    troopBudgetMultiplier,
	}
}
