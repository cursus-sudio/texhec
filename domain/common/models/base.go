package models

import "github.com/ogiusek/ioc"

type ModelBase struct {
	Id ModelId
}

func (base *ModelBase) Valid(c ioc.Dic) []error {
	// TODO
	return nil
}

func NewBase(c ioc.Dic) ModelBase {
	return ModelBase{
		// TODO
	}
}
