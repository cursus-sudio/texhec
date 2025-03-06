package models

// identifier

type ModelId string

func (id ModelId) String() string {
	return string(id)
}
