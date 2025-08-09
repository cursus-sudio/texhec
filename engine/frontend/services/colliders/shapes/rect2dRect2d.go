package shapes

import (
	"frontend/services/colliders"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

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

func (s *collidersService) rect2DRect2DHandler(s1 Rect2D, s2 Rect2D) colliders.Collision {
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
