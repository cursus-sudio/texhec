package errdesc

type ErrorWithPath interface {
	error

	// returns only message
	Message() string

	// returns only path
	Path() string

	// adds property name to path
	// parse property name to which you want to display
	//
	// note: this method modifies and returns modified version for convenience
	Property(string) ErrorWithPath
}

func ErrPath(err error) ErrorWithPath {
	if errWithPath, ok := err.(ErrorWithPath); ok {
		return errWithPath
	}
	return WithDescription(err)
}

func Path(err error) error {
	return err
}
