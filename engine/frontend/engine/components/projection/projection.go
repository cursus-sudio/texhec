package projection

import (
	"errors"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	ErrWorldShouldHaveOneProjection error = errors.New("world should have one projection")
)

type Projection struct {
	Projection mgl32.Mat4
}

func NewProjection(projection mgl32.Mat4) Projection {
	return Projection{projection}
}
