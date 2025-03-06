package models

type ModelDescription interface {
	Name() string
	Description() string
	Photos() []string
}

type ModelDescriptionSetter interface {
}

// type ModelDescription struct {
// 	name        string
// 	description string
// 	photoId     string
// }
