package codec

import "errors"

var (
	ErrTypeIsNotRegistered     error = errors.New("type is not registered")
	ErrTypeIsAlreadyRegistered error = errors.New("type is already registered")

	ErrInvalidBytes     error = errors.New("invalid bytes")
	ErrCannotEncodeType error = errors.New("cannot encode type")
)
