package blueprints

import (
	"domain/common/models"
	"domain/vo"
)

type SkillInGame struct {
	models.ModelBase
	models.ModelDescription
	// Action  skills.SkillAction
	// UseRule skills.SkillUseRule
	Cost vo.TroopBudget
}
