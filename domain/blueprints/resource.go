package blueprints

import (
	"domain/common/models"

	"github.com/ogiusek/ioc"
)

type Resource struct {
	models.ModelBase
	models.ModelDescription
	MinDeposit, MaxDeposit uint
}

func NewResource(c ioc.Dic, desc models.ModelDescription, minDeposit, maxDeposit uint) Resource {
	return Resource{
		ModelBase:        models.NewBase(c),
		ModelDescription: desc,
		MinDeposit:       minDeposit,
		MaxDeposit:       maxDeposit,
	}
}
