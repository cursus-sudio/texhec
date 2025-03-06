package gamedefinitions

import (
	"domain/common/models"
	"domain/states/kernel"
)

type Fraction interface {
	models.ModelBase
	models.ModelDescription
	InitialResearch() kernel.ResearchDefinition

	// models.ModelBase
	// models.ModelDescription
	// spawn           Biom
	// brain           Building
	// initialResearch kernel.ResearchDefinition
}
