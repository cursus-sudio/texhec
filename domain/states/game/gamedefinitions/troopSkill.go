package gamedefinitions

import "domain/common/models"

// type TroopSkillTarget uint

// const (
// 	TroopSkillTargetBuilding         TroopSkillTarget = iota
// 	TroopSkillTargetNothing          TroopSkillTarget = iota
// 	TroopSkillTargetEnemyTroop       TroopSkillTarget = iota
// 	TroopSkillTargetFriendlyTroop    TroopSkillTarget = iota
// 	TroopSkillTarget2X2EmptySpace    TroopSkillTarget = iota
// 	TroopSkillTarget2X2MountainSpace TroopSkillTarget = iota
// )

type TroopSkill interface {
	models.ModelBase
	models.ModelDescription

	HasTarget() bool
	SkillId() string
}
