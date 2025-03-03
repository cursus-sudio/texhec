package blueprints

import (
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

type Research struct {
	models.ModelBase
	models.ModelDescription
	Cost                    vo.Budget
	RequiredResearchUnlocks uint // 1 by default for more advanced things like behemot can be 2 or 3
	UnlocksResearch         []*Research
	UnlocksTroops           []*Troop
	UnlocksBuildings        []*Building
}

func NewResearch(
	c ioc.Dic,
	desc models.ModelDescription,
	cost vo.Budget,
	requiredResearchUnlocks uint,
	unlocksResearch []*Research,
	unlocksTroops []*Troop,
	unlocksBuildings []*Building,
) Research {
	return Research{
		ModelBase:               models.NewBase(c),
		ModelDescription:        desc,
		Cost:                    cost,
		RequiredResearchUnlocks: requiredResearchUnlocks,
		UnlocksResearch:         unlocksResearch,
		UnlocksTroops:           unlocksTroops,
		UnlocksBuildings:        unlocksBuildings,
	}
}
