package saves

import (
	"errors"
)

type SaveMetaFactory interface {
	// can return Err400
	New(SaveName) SaveMeta
}

// func NewSaveMeta(name SaveName) SaveMeta {}

var (
	ErrInvalidQuery    error = errors.New("save storage: invalid query")
	ErrSaveDoNotExists error = errors.New("save storage: save with this id do not exists")
)

type Saves interface {
	// panics when query is invalid
	ListSaves(ListSavesQuery) ([]SaveMeta, error)
	SavesPages(ListSavesQuery) (int, error)

	// can return:
	// - Err404
	// - Err422
	// - Err503
	Load(SaveId) error

	// overwrites save if it already exists
	// can return:
	// - Err404
	// - Err503
	Save(SaveId) error

	// can return:
	// - Err503
	Delete(SaveId) error

	// if save already exists it overwrites it
	// can return:
	// - Err503
	NewSave(SaveMeta) error
}
