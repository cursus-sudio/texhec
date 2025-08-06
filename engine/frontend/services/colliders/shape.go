package colliders

import (
	"frontend/engine/components/transform"

	"github.com/go-gl/mathgl/mgl32"
)

type Shape interface {
	Position() mgl32.Vec3
	Apply(transform.Transform) Shape
}
