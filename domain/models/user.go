package models

import (
	"domain/common/models"
	"domain/vo"

	"github.com/ogiusek/ioc"
)

type User struct {
	models.ModelBase
	models.ModelDescription
	Password vo.Hash
}

func NewUser(c ioc.Dic, desc models.ModelDescription, password vo.Hash) User {
	return User{
		ModelBase:        models.NewBase(c),
		ModelDescription: desc,
		Password:         password,
	}
}
