package blueprints

import (
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
	"github.com/ogiusek/null"
)

type Building struct {
	models.ModelBase
	models.ModelDescription
	Cost             vo.Budget
	OnResource       null.Nullable[*Resource]
	HP               vo.HP
	Income           vo.Budget
	Maintanace       vo.Budget
	ResearchCapacity vo.Budget
	Recruits         []*Troop
}

func NewBuilding(
	c ioc.Dic,
	description models.ModelDescription,
	cost vo.Budget,
	onResource null.Nullable[*Resource],
	hp vo.HP,
	income vo.Budget,
	researchCapacity vo.Budget,
	recruits []*Troop,
) Building {
	return Building{
		ModelBase:        models.NewBase(c),
		ModelDescription: description,
		Cost:             cost,
		OnResource:       onResource,
		HP:               hp,
		Income:           income,
		ResearchCapacity: researchCapacity,
		Recruits:         recruits,
	}
}

func (build *Building) CanBuildOn(resource null.Nullable[*Resource]) bool {
	return build.OnResource.Ok == resource.Ok && (!resource.Ok || build.OnResource.Val.Id == resource.Val.Id)
}
