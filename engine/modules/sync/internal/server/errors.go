package server

import "errors"

var (
	ErrRecordingDidntStartProperly error = errors.New("recording didn't start properly")
)
