package errdesc

import "fmt"

type errorWithDesc struct {
	message string
	path    string
}

func (err *errorWithDesc) Error() string {
	return fmt.Sprintf("`%s` -> %s", err.path, err.message)
}

func (err *errorWithDesc) Message() string {
	return err.message
}

func (err *errorWithDesc) Path() string {
	return err.path
}

func (err *errorWithDesc) Property(parent string) ErrorWithPath {
	if err.path == "" {
		err.path = parent
	} else {
		err.path = fmt.Sprintf("%s.%s", parent, err.path)
	}
	return err
}

type errDesc interface {
	ErrorWithPath
}

func WithDescription(err error) errDesc {
	return &errorWithDesc{
		message: err.Error(),
		path:    "",
	}
}
