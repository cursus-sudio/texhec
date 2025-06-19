package files

// path is just generally unique file id
type Path string

func NewPath(path string) Path {
	return Path(path)
}

func (path *Path) String() string {
	return string(*path)
}

// guid for example
type PathGenerator interface {
	New() Path
}

// this is blob repository
type FileStorage interface {
	// creates file if it does not exist
	// returned errors:
	// - Err400
	// - Err503
	EnsureExists(Path) error

	// returns true if file with this Path exists
	// returned errors:
	// - Err503
	Exists(Path) (bool, error)

	// returned errors:
	// - Err404
	// - Err503
	Read(Path) ([]byte, error)

	// if file do not exists creates it
	// returned errors:
	// - Err503
	OverWrite(Path, []byte) error

	// if file do not exists creates it
	// returned errors:
	// - Err503
	Write(Path, []byte) error

	Delete(Path) error
}
