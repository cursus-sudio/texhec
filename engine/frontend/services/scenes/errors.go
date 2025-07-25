package scenes

import (
	"errors"
)

var (
	ErrSceneAlreadyExists error = errors.New("scene already exists")
	ErrNoActiveScene      error = errors.New("no active scene")
	ErrSceneDoNotExists   error = errors.New("scene do not exists")
)
