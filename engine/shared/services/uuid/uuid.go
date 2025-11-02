package uuid

import (
	"github.com/google/uuid"
	"github.com/ogiusek/ioc/v2"
)

// interface

type UUID struct {
	val uuid.UUID
}

func newUUID(val uuid.UUID) UUID {
	return UUID{val: val}
}

func (uuid UUID) String() string { return uuid.val.String() }
func (uuid UUID) Bytes() []byte  { return uuid.val[:] }

type Factory interface {
	NewUUID() UUID
}

// impl

type factory struct{}

func (factory *factory) NewUUID() UUID {
	return newUUID(uuid.New())
}

// pkg

type pkg struct{}

func Package() ioc.Pkg {
	return pkg{}
}

func (pkg pkg) Register(b ioc.Builder) {
	ioc.RegisterSingleton(b, func(c ioc.Dic) Factory { return &factory{} })
}
