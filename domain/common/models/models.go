package models

import "github.com/ogiusek/ioc"

type ModelDescription struct {
	Name        string
	Description string
	PhotoId     string
}

func (*ModelDescription) Valid(ioc.Dic) []error {
	// TODO
	return nil
}
