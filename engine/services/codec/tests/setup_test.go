package test

import (
	"engine/services/codec"
	"engine/services/logger"

	"github.com/ogiusek/ioc/v2"
)

type Type struct {
	Value int
}

type setup struct {
	codec codec.Codec
}

func NewSetup() setup {
	b := ioc.NewBuilder()

	for _, pkg := range []ioc.Pkg{
		codec.Package(),
		logger.Package(true, func(c ioc.Dic, message string) { print(message) }),
	} {
		pkg.Register(b)
	}

	ioc.WrapService(b, ioc.DefaultOrder, func(c ioc.Dic, b codec.Builder) codec.Builder {
		return b.Register(Type{})
	})

	c := b.Build()
	return setup{ioc.Get[codec.Codec](c)}
}
