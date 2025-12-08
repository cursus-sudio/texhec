package client

import "errors"

var (
	ErrInvalidPrediction    error = errors.New("client had invalid prediction")
	ErrExceededPredictions  error = errors.New("predictions count exceed maximal amount")
	ErrHasMoreThanOneServer error = errors.New("has more than one server")
)
