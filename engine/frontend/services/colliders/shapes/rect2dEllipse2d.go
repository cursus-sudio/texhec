package shapes

import (
	"frontend/services/colliders"

	"github.com/go-gl/mathgl/mgl32"
)

func (s *collidersService) rect2DEllipse2DHandler(s1 Rect2D, s2 Ellipsoid2D) colliders.Collision {
	rect := s1
	circle := s2

	inverseRot := rect.Rotation.Inverse()
	rectCenter := rect.Pos
	circleCenter := circle.Pos
	circleRadius := circle.R()

	v := circleCenter.Sub(rectCenter)

	vRotated := inverseRot.Rotate(v)

	halfWidth := rect.Size[0] / 2
	halfHeight := rect.Size[1] / 2
	closestPointLocal := mgl32.Vec3{
		max(-halfWidth, min(halfWidth, vRotated.X())),
		max(-halfHeight, min(halfHeight, vRotated.Y())),
		0,
	}

	distVecLocal := vRotated.Sub(closestPointLocal)

	if distVecLocal.LenSqr() == 0 {
		overlapX := halfWidth - mgl32.Abs(vRotated.X())
		overlapY := halfHeight - mgl32.Abs(vRotated.Y())

		var normal mgl32.Vec3
		var penetrationDepth float32

		if overlapX < overlapY {
			normal = rect.Rotation.Rotate(mgl32.Vec3{1, 0, 0})
			if vRotated.X() < 0 {
				normal = normal.Mul(-1)
			}
			penetrationDepth = circleRadius + overlapX
		} else {
			normal = rect.Rotation.Rotate(mgl32.Vec3{0, 1, 0})
			if vRotated.Y() < 0 {
				normal = normal.Mul(-1)
			}
			penetrationDepth = circleRadius + overlapY
		}

		contactPointRect := circleCenter.Sub(normal.Mul(overlapX))
		contactPointCircle := circleCenter.Sub(normal.Mul(circleRadius))

		intersection := colliders.NewIntersection(
			contactPointRect,
			contactPointCircle,
			normal,
			penetrationDepth,
		)
		return colliders.NewCollision(intersection)
	}

	closestPointWorld := rectCenter.Add(rect.Rotation.Rotate(closestPointLocal))

	distVecWorld := circleCenter.Sub(closestPointWorld)
	distSq := distVecWorld.LenSqr()
	radiusSq := pow2(circleRadius)

	if distSq > radiusSq {
		return nil
	}

	dist := sqrt2(distSq)
	normal := distVecWorld.Mul(1 / dist)
	penetrationDepth := circleRadius - dist

	contactPointRect := closestPointWorld
	contactPointCircle := circleCenter.Sub(normal.Mul(circleRadius))

	intersection := colliders.NewIntersection(
		contactPointRect,
		contactPointCircle,
		normal,
		penetrationDepth,
	)

	return colliders.NewCollision(intersection)
}
