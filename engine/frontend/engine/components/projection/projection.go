package projection

import (
	"errors"
)

var (
	ErrWorldShouldHaveOneProjection error = errors.New("world should have one projection")
)

// type Projection struct {
// 	Projection mgl32.Mat4
// }
//
// func NewProjection(projection mgl32.Mat4) Projection {
// 	return Projection{Projection: projection}
// }

type Perspecrive struct {
	FovY        float32
	AspectRatio float32
	Near, Far   float32
}

func NewPerspective(fovY float32, aspectRatio float32, near, far float32) Perspecrive {
	return Perspecrive{FovY: fovY, AspectRatio: aspectRatio, Near: near, Far: far}
}

type Ortho struct {
	Width, Height float32
	Near, Far     float32
}

func NewOrtho(w, h, near, far float32) Ortho {
	return Ortho{Width: w, Height: h, Near: near, Far: far}
}
