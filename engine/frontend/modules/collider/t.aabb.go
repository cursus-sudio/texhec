package collider

import (
	"github.com/go-gl/mathgl/mgl32"
)

type AABB struct {
	Min, Max mgl32.Vec3
}

func NewAABB(min, max mgl32.Vec3) AABB {
	return AABB{Min: min, Max: max}
}
