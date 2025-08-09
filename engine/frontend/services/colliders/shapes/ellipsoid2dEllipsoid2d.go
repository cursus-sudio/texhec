package shapes

import (
	"frontend/services/colliders"

	"github.com/go-gl/mathgl/mgl32"
)

func (s *collidersService) ellipsoid2DEllipsoid2DHandler(s1 Ellipsoid2D, s2 Ellipsoid2D) colliders.Collision {
	if s1.Pos == s2.Pos {
		return colliders.NewCollision()
	}
	distSq := pow2(s1.Pos.X()-s2.Pos.X()) + pow2(s1.Pos.Y()-s2.Pos.Y())
	radiusSumSq := pow2(s1.R() + s2.R())
	if distSq > radiusSumSq {
		return nil
	}
	dist := sqrt2(distSq)
	penetrationDepth := (s1.R() + s2.R()) - float32(dist)
	normal := mgl32.Vec3{s2.Pos[0], s2.Pos[1]}.
		Sub(s1.Pos).
		Mul(1 / dist)

	contactPoint1 := mgl32.Vec3{s1.Pos[0], s1.Pos[1]}.
		Add(normal.Mul(s1.R()))

	contactPoint2 := mgl32.Vec3{s2.Pos[0], s2.Pos[1]}.
		Add(normal.Mul(-s2.R()))

	intersection := colliders.NewIntersection(
		contactPoint1,
		contactPoint2,
		normal,
		penetrationDepth,
	)

	return colliders.NewCollision(intersection)
}
