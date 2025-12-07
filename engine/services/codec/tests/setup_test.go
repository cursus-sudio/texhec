package test

import "engine/services/codec"

type Type struct {
	Value int
}

type setup struct {
	codec codec.Codec
}

func NewSetup() setup {
	b := codec.NewBuilder()
	b.Register(Type{})
	c := b.Build()
	return setup{c}
}
