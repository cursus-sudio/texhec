package models

type ModelBase interface {
	Id() ModelId
}

type ModelBaseSetter interface{}
