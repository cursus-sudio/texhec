package gamedefinitions

import (
	"domain/common/models"
	"domain/states/kernel/vo"
)

type Building interface {
	models.ModelBase
	models.ModelDescription
	// models.ModelBase
	// models.ModelDescription
	// maxHp    vo.Hp
	// metadata gameobject.Metadata
	MaxHp() vo.Hp
}
