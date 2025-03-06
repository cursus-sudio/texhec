package general

import (
	"domain/common/models"
	uservo "domain/states/user/vo"

	"github.com/ogiusek/ioc"
)

type User struct {
	models.ModelBase
	models.ModelDescription
	Password uservo.Hash
}

func NewUser(c ioc.Dic, base models.ModelBase, desc models.ModelDescription, password uservo.Hash) User {
	return User{
		ModelBase:        base,
		ModelDescription: desc,
		Password:         password,
	}
}
