package gameobjects

import (
	"domain/common/models"
	"domain/states/kernel"
	"domain/states/kernel/vo"

	"github.com/ogiusek/ioc"
)

type Research struct {
	models.ModelBase
	definition kernel.ResearchDefinition
	progress   vo.Budget
}

func NewResearch(c ioc.Dic, base models.ModelBase, definition kernel.ResearchDefinition, progress vo.Budget) Research {
	return Research{
		ModelBase:  base,
		definition: definition,
		progress:   progress,
	}
}

func (research *Research) Definition() kernel.ResearchDefinition {
	return research.definition
}

func (research *Research) Progress() vo.Budget {
	return research.progress
}

func (research *Research) IsComplete() bool {
	return research.definition.Cost() == research.progress
}

// returns overflow vo.Budget
func (r *Research) Research(budget vo.Budget) vo.Budget {
	max := r.definition.Cost()
	var overflow vo.Budget
	r.progress, overflow = max.MaxSum(r.progress, budget)
	return overflow
}
