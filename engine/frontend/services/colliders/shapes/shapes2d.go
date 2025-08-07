package shapes

import (
	"frontend/engine/components/transform"
	"frontend/services/colliders"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// circle
// rectangle
// stadium (not needed for now)

type ellipsoid2D struct {
	transform.Transform
}

type Ellipsoid2D struct{ ellipsoid2D }

func (s ellipsoid2D) R() float32 { return max(s.Size[0], s.Size[1]) }

func (s ellipsoid2D) Apply(t transform.Transform) colliders.Shape {
	return Ellipsoid2D{ellipsoid2D{s.Transform.Merge(t)}}
}

func (s ellipsoid2D) Position() mgl32.Vec3 { return s.Pos }

// We currently implement only circles but store it as ellipse for further development
// but currently treat ellipse2D as circle
// func NewEllipsoid2D(t transform.Transform) colliders.Shape {
// 	return Ellipsoid2D{ellipsoid2D{t}}
// }

func NewCircle2D(pos mgl32.Vec3, r float32) colliders.Shape {
	return Ellipsoid2D{ellipsoid2D{
		transform.NewTransform().
			SetPos(pos).
			SetSize(mgl32.Vec3{r, r, r}),
	}}
}

func ellipsoid2DEllipsoid2DHandler(s1 Ellipsoid2D, s2 Ellipsoid2D) colliders.Collision {
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

//

type rect2D struct {
	transform.Transform
}

type Rect2D struct{ rect2D }

func (s rect2D) Apply(t transform.Transform) colliders.Shape {
	return Rect2D{rect2D{s.Transform.Merge(t)}}
}

func (s rect2D) Position() mgl32.Vec3 { return s.Pos }

func NewRect2D(t transform.Transform) Rect2D {
	return Rect2D{rect2D{t}}
}

func project(vertices []mgl32.Vec3, axis mgl32.Vec3) (minProjection, maxProjection float32) {
	minProjection = vertices[0].Dot(axis)
	maxProjection = minProjection
	for i := 1; i < len(vertices); i++ {
		p := vertices[i].Dot(axis)
		if p < minProjection {
			minProjection = p
		} else if p > maxProjection {
			maxProjection = p
		}
	}
	return minProjection, maxProjection
}

func getOverlap(min1, max1, min2, max2 float32) (isOverlapping bool, penetration float32) {
	if max1 < min2 || max2 < min1 {
		return false, 0
	}
	if max1 > max2 {
		penetration = max2 - min1
	} else {
		penetration = max1 - min2
	}
	return true, penetration
}

func getRectVertices(r Rect2D) []mgl32.Vec3 {
	hw := r.Size[0] / 2
	hh := r.Size[1] / 2

	localVertices := []mgl32.Vec3{
		{-hw, -hh, 0},
		{+hw, -hh, 0},
		{+hw, +hh, 0},
		{-hw, +hh, 0},
	}

	var vertices []mgl32.Vec3
	rotMat := r.Rotation.Mat4()
	for _, v := range localVertices {
		rotatedVertex := rotMat.Mul4x1(v.Vec4(1)).Vec3()
		translatedVertex := rotatedVertex.Add(r.Pos)
		vertices = append(vertices, translatedVertex)
	}

	return vertices
}

func rect2DRect2DHandler(s1 Rect2D, s2 Rect2D) colliders.Collision {
	if s1.Pos == s2.Pos {
		return colliders.NewCollision()
	}
	vertices1 := getRectVertices(s1)
	vertices2 := getRectVertices(s2)

	axes := make([]mgl32.Vec3, 4)

	axes[0] = vertices1[1].Sub(vertices1[0]).Normalize()
	axes[1] = vertices1[2].Sub(vertices1[1]).Normalize()

	axes[2] = vertices2[1].Sub(vertices2[0]).Normalize()
	axes[3] = vertices2[2].Sub(vertices2[1]).Normalize()

	var minPenetration float32
	var collisionNormal mgl32.Vec3
	var firstOverlap = true

	for _, axis := range axes {
		min1, max1 := project(vertices1, axis)
		min2, max2 := project(vertices2, axis)

		isOverlapping, penetration := getOverlap(min1, max1, min2, max2)
		if !isOverlapping {
			return nil
		}

		if firstOverlap || penetration < minPenetration {
			minPenetration = penetration
			collisionNormal = axis
			firstOverlap = false
		}
	}

	center1 := s1.Pos
	center2 := s2.Pos
	v := center2.Sub(center1)
	if v.Dot(collisionNormal) < 0 {
		collisionNormal = collisionNormal.Mul(-1)
	}

	maxDot1 := float32(math.MaxFloat32) * -1
	var contactPoint1 mgl32.Vec3
	for _, v := range vertices1 {
		dot := v.Dot(collisionNormal)
		if dot > maxDot1 {
			maxDot1 = dot
			contactPoint1 = v
		}
	}

	maxDot2 := float32(math.MaxFloat32) * -1
	var contactPoint2 mgl32.Vec3
	for _, v := range vertices2 {
		dot := v.Dot(collisionNormal.Mul(-1))
		if dot > maxDot2 {
			maxDot2 = dot
			contactPoint2 = v
		}
	}

	intersection := colliders.NewIntersection(
		contactPoint1,
		contactPoint2,
		collisionNormal,
		minPenetration,
	)

	return colliders.NewCollision(intersection)
}

func rect2DEllipse2DHandler(s1 Rect2D, s2 Ellipsoid2D) colliders.Collision {
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

func rect2DRayHandler(s1 Rect2D, s2 Ray) colliders.Collision {
	rect := s1
	rayOrigin := s2.Pos
	rayDirection := s2.Rotation.Rotate(mgl32.Vec3{0, 1, 0}).Normalize()

	if rayDirection.LenSqr() < 1e-6 {
		return nil
	}

	inverseRectRotation := rect.Rotation.Inverse()
	localRayOrigin := inverseRectRotation.Rotate(rayOrigin.Sub(rect.Pos))
	localRayDirection := inverseRectRotation.Rotate(rayDirection)

	halfWidth := rect.Size[0] / 2
	halfHeight := rect.Size[1] / 2
	minBounds := mgl32.Vec3{-halfWidth, -halfHeight, 0}
	maxBounds := mgl32.Vec3{halfWidth, halfHeight, 0}

	tEntry := float32(math.Inf(-1))
	tExit := float32(math.Inf(1))

	var normalEntry mgl32.Vec3

	for i := 0; i < 2; i++ {
		if mgl32.Abs(localRayDirection[i]) < 1e-6 {
			if localRayOrigin[i] < minBounds[i] || localRayOrigin[i] > maxBounds[i] {
				return nil
			}
		} else {
			invDir := 1.0 / localRayDirection[i]
			tNear := (minBounds[i] - localRayOrigin[i]) * invDir
			tFar := (maxBounds[i] - localRayOrigin[i]) * invDir

			if tNear > tFar {
				tNear, tFar = tFar, tNear
			}

			if tNear > tEntry {
				tEntry = tNear
				normalEntry = mgl32.Vec3{}
				if invDir < 0 {
					normalEntry[i] = 1.0
				} else {
					normalEntry[i] = -1.0
				}
			}
			tExit = min(tExit, tFar)

			if tEntry > tExit {
				return nil
			}
		}
	}

	if tExit < 0 {
		return nil
	}

	var hitT float32
	var finalLocalNormal mgl32.Vec3

	if tEntry < 0 {
		hitT = tExit
		finalLocalNormal = normalEntry.Mul(-1)
	} else {
		hitT = tEntry
		finalLocalNormal = normalEntry
	}

	localContactPoint := localRayOrigin.Add(localRayDirection.Mul(hitT))

	worldContactPointRect := rect.Pos.Add(rect.Rotation.Rotate(localContactPoint))
	worldContactPointRay := rayOrigin.Add(rayDirection.Mul(hitT))

	worldNormal := rect.Rotation.Rotate(finalLocalNormal).Normalize()

	penetrationDepth := float32(0)

	intersection := colliders.NewIntersection(
		worldContactPointRect,
		worldContactPointRay,
		worldNormal,
		penetrationDepth,
	)

	return colliders.NewCollision(intersection)
}

func ellipse2DRayHandler(s1 Ellipsoid2D, s2 Ray) colliders.Collision {
	circle := s1
	ray := s2

	circlePos := circle.Pos
	circleRadius := circle.R()

	rayPos := ray.Pos
	rayDir := ray.Rotation.Rotate(mgl32.Vec3{0, 1, 0}).Normalize()

	// Remove the Z-coordinate checks
	// The original code checked for Z-plane alignment, which is not desired for 2D collision with Z as an index.
	// if mgl32.Abs(rayDir.Z()) < 1e-6 {
	//     if mgl32.Abs(rayPos.Z()-circlePos.Z()) > 1e-6 {
	//         return nil
	//     }
	// } else {
	//     tPlane := (circlePos.Z() - rayPos.Z()) / rayDir.Z()
	//     if tPlane < 0 {
	//         return nil
	//     }
	// }

	circlePos2D := mgl32.Vec2{circlePos.X(), circlePos.Y()}
	rayPos2D := mgl32.Vec2{rayPos.X(), rayPos.Y()}
	rayDir2D := mgl32.Vec2{rayDir.X(), rayDir.Y()}

	m2D := rayPos2D.Sub(circlePos2D)

	if rayDir2D.LenSqr() < 1e-6 {
		if m2D.LenSqr() <= circleRadius*circleRadius {
			normal2D := m2D.Normalize()
			contactPoint3D := mgl32.Vec3{rayPos2D.X(), rayPos2D.Y(), circlePos.Z()}
			normal3D := mgl32.Vec3{normal2D.X(), normal2D.Y(), 0}

			intersection := colliders.NewIntersection(
				contactPoint3D,
				contactPoint3D,
				normal3D,
				0,
			)
			return colliders.NewCollision(intersection)
		}
		return nil
	}
	a := rayDir2D.Dot(rayDir2D)
	b := 2 * m2D.Dot(rayDir2D)
	c := m2D.Dot(m2D) - circleRadius*circleRadius

	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return nil
	}

	sqrtDiscriminant := float32(math.Sqrt(float64(discriminant)))
	t1 := (-b - sqrtDiscriminant) / (2 * a)
	t2 := (-b + sqrtDiscriminant) / (2 * a)

	var finalT float32 = -1.0

	if t1 >= 0 {
		finalT = t1
	}

	if t2 >= 0 {
		if finalT < 0 || t2 < finalT {
			finalT = t2
		}
	}

	if finalT < 0 {
		return nil
	}

	intersectionPoint2D := rayPos2D.Add(rayDir2D.Mul(finalT))
	intersectionPoint3D := mgl32.Vec3{intersectionPoint2D.X(), intersectionPoint2D.Y(), circlePos.Z()}

	normal2D := intersectionPoint2D.Sub(circlePos2D).Normalize()
	normal3D := mgl32.Vec3{normal2D.X(), normal2D.Y(), 0}

	intersection := colliders.NewIntersection(
		intersectionPoint3D,
		intersectionPoint3D,
		normal3D,
		0,
	)

	return colliders.NewCollision(intersection)
}
