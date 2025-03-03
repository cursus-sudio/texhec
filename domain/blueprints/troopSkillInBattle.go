package blueprints

import (
	"domain/common/errdesc"
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

// TODO
// REDO

// SERVICE
type SkillErrors struct {
	UsedInDoesNotMatch error
}

type SkillInBattle struct {
	models.ModelBase
	models.ModelDescription
	Cost vo.TroopBudget
}

func (skill *SkillInBattle) Valid(c ioc.Dic) []error {
	// errors := ioc.Get[SkillErrors](c)
	var errs []error
	errs = append(errs, skill.ModelBase.Valid(c)...)
	errs = append(errs, skill.ModelDescription.Valid(c)...)

	for _, err := range skill.Cost.Valid(c) {
		errs = append(errs, errdesc.ErrPath(err).Property("cost"))
	}
	return errs
}

func NewSkillInBattle(
	c ioc.Dic,
	desc models.ModelDescription,
	cost vo.TroopBudget,
) (*SkillInBattle, []error) {
	skill := &SkillInBattle{
		ModelBase:        models.NewBase(c),
		ModelDescription: desc,
		Cost:             cost,
	}
	return skill, skill.Valid(c)
}
