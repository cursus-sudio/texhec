package scopes

import "errors"

var (
	ErrAlreadyCleanedUp error = errors.New("scope clean up already cleaned up")
)
