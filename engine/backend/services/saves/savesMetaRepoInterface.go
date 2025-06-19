package saves

import (
	"time"
)

type SaveName string

type SaveMeta struct {
	Id SaveId

	Created      time.Time
	LastModified time.Time

	Name SaveName
	// optional to add:
	// - preview
}

func NewSaveMeta(
	id SaveId,
	created time.Time,
	name SaveName,
) SaveMeta {
	return SaveMeta{
		Id:           id,
		Created:      created,
		LastModified: created,
		Name:         name,
	}
}

type SavesMetaRepo interface {
	// Create
	// Update
	// can return:
	// - Err503
	Upsert(SaveMeta) error

	// Delete
	// can return:
	// - Err503
	Delete(SaveId) error

	// Read
	// can return:
	// - Err503
	ListSaves(ListSavesQuery) ([]SaveMeta, error)

	// can return:
	// - Err503
	SavesPages(ListSavesQuery) (int, error)

	// can return:
	// - Err503
	// - Err404
	GetById(SaveId) (SaveMeta, error)
}
