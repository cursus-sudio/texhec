package game

import "github.com/ogiusek/ioc"

type HubName string

func (*HubName) Valid(c ioc.Dic) []error {
	// TODO
	return nil
}
