package gamedefinitions

import "domain/common/models"

type Resource interface {
	models.ModelBase
	models.ModelDescription
	// models.ModelBase
	// models.ModelDescription
	// minDeposit, maxDeposit uint
}
