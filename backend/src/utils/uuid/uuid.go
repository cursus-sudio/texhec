package uuid

import (
	"github.com/google/uuid"
	"github.com/ogiusek/ioc"
)

// interface

type UUID string

func NewUUID(val string) UUID {
	return UUID(val)
}

func (uuid UUID) String() string {
	return string(uuid)
}

type Factory interface {
	NewUUID() UUID
}

// impl

type factory struct{}

func (factory *factory) NewUUID() UUID {
	return NewUUID(uuid.New().String())
}

// pkg

type Pkg struct{}

func Package() Pkg {
	return Pkg{}
}

func (pkg Pkg) Register(c ioc.Dic) {
	ioc.RegisterSingleton(c, func(c ioc.Dic) Factory { return &factory{} })
}
