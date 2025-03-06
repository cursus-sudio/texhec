package kernel

import (
	"domain/common/models"
	"domain/states/kernel/vo"

	"github.com/ogiusek/ioc"
)

// SERVICE
type ResearchErrors interface {
	RequiredUnlocksMustBePositive() error
}

type ResearchDefinition struct {
	models.ModelBase
	models.ModelDescription
	cost            vo.Budget
	requiredUnlocks uint // 1 by default for more advanced things like behemot can be 2 or 3
	unlocksResearch []ResearchDefinition
}

func NewResearchDefinition(
	c ioc.Dic,
	base models.ModelBase,
	desc models.ModelDescription,
	cost vo.Budget,
	unlocksResearch []ResearchDefinition,
) ResearchDefinition {
	return ResearchDefinition{
		ModelBase:        base,
		ModelDescription: desc,
		cost:             cost,
		requiredUnlocks:  1,
		unlocksResearch:  unlocksResearch,
	}
}

func (r *ResearchDefinition) SetRequiredResearches(c ioc.Dic, requiredUnlocks uint) error {
	researchErrors := ioc.Get[ResearchErrors](c)
	if requiredUnlocks == 0 {
		return researchErrors.RequiredUnlocksMustBePositive()
	}
	r.requiredUnlocks = requiredUnlocks
	return nil
}

func (r *ResearchDefinition) Cost() vo.Budget {
	return r.cost
}

func (r *ResearchDefinition) NextResearches() []ResearchDefinition {
	return r.unlocksResearch
}
