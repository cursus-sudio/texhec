package gamedefinitions

import (
	"domain/common/gameobject"
	"domain/common/models"
	"domain/states/kernel/vo"
)

// SERVICE
type TroopErrors interface {
	SkillDoNotExist(TroopSkill) error
	SkillAlreadyExists(TroopSkill) error
}

type Troop interface {
	models.ModelBase
	models.ModelDescription
	// models.ModelBase
	// kernel.TroopDefinition
	// metadata gameobject.Metadata
	// Skills   map[models.ModelId]*TroopSkill

	// troops should be able to:
	// - move
	// - attack buildings
	// - build
	// - see (remove fog of war)
	// - kill itself
	// - heal itself
	// - many others
	TroopBudget() vo.TroopBudget
	MaxAmount() uint

	Metadata() gameobject.Metadata
	SetMetadata(m gameobject.Metadata)
	GetSkills() map[models.ModelId]TroopSkill

	// AddSkill(skill TroopSkill) error
	// RemoveSkill(skill TroopSkill) error
}
