package internal

import (
	"engine/modules/uuid"
	uuidSource "github.com/google/uuid"
)

// impl

type factory struct{}

func (factory *factory) NewUUID() uuid.UUID {
	return uuid.UUID(uuidSource.New())
}

func NewFactory() uuid.Factory {
	return &factory{}
}
