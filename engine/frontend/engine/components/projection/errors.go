package projection

import (
	"errors"
)

var (
	ErrWorldShouldHaveOneProjection error = errors.New("world should have one projection")
)
