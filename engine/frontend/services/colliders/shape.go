package colliders

import (
	"frontend/engine/components/transform"
)

type Shape interface {
	// Position() mgl32.Vec3
	Apply(transform.Transform) Shape
}
