package saves

type SaveId string

func NewSaveId(id string) SaveId {
	return SaveId(id)
}

func (id *SaveId) String() string {
	return string(*id)
}

// this storage is responsible for storing saves and loading them.
// its not responsible for querying them.
// this is just overlay for blob repository and loading them.
type SavesStorage interface {
	// can return:
	// - httperrors.Err400
	// - httperrors.Err404
	// - httperrors.Err422 (returned when stored save is invalid)
	// - httperrors.Err503
	LoadSave(SaveId) error

	// can return:
	// -httperrors.Err503
	Save(SaveId) error

	// can return:
	// -httperrors.Err503
	Delete(SaveId) error
}
