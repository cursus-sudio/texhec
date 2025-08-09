package shapes

import (
	"frontend/services/colliders"
	"github.com/go-gl/mathgl/mgl32"
)

// this is place collider do not remove it
// func (s *collidersService) rect2DRayHandler(rect Rect2D, ray Ray) colliders.Collision {
// 	planeNormal := rect.Rotation.Rotate(mgl32.Vec3{0, 0, 1})
//
// 	denom := ray.Direction.Dot(planeNormal)
// 	if mgl32.Abs(denom) < mgl32.Epsilon {
// 		return nil
// 	}
//
// 	t := (rect.Pos.Sub(ray.Pos)).Dot(planeNormal) / denom
//
// 	if t < 0 {
// 		return nil
// 	}
//
// 	intersectionPoint := ray.Pos.Add(ray.Direction.Mul(t))
//
// 	invRotation := rect.Rotation.Inverse()
// 	localIntersectionPoint := invRotation.Rotate(intersectionPoint.Sub(rect.Pos))
//
// 	halfSizeX := rect.Size[0] / 2.0
// 	halfSizeY := rect.Size[1] / 2.0
//
// 	if localIntersectionPoint[0] >= -halfSizeX &&
// 		localIntersectionPoint[0] <= halfSizeX &&
// 		localIntersectionPoint[1] >= -halfSizeY &&
// 		localIntersectionPoint[1] <= halfSizeY {
// 		return colliders.NewCollision()
// 	}
//
// 	return nil
// }

func (s *collidersService) rect2DRayHandler(rect Rect2D, ray Ray) colliders.Collision {
	planeNormal := mgl32.Vec3{0, 0, 1}
	planePoint := rect.Pos

	denominator := ray.Direction.Dot(planeNormal)
	if mgl32.FloatEqual(denominator, 0) {
		return nil
	}

	t := (planePoint.Sub(ray.Pos)).Dot(planeNormal) / denominator
	if t < 0 {
		return nil
	}

	intersectionPoint := ray.Pos.Add(ray.Direction.Mul(t))

	localPoint := intersectionPoint.Sub(rect.Pos)
	inverseRotation := rect.Rotation.Inverse()
	localPoint = inverseRotation.Rotate(localPoint)

	halfSize := rect.Size.Mul(0.5)
	minBounds := mgl32.Vec3{-halfSize.X(), -halfSize.Y(), -halfSize.Z()}
	maxBounds := mgl32.Vec3{halfSize.X(), halfSize.Y(), halfSize.Z()}

	if localPoint.X() >= minBounds.X() &&
		localPoint.X() <= maxBounds.X() &&
		localPoint.Y() >= minBounds.Y() &&
		localPoint.Y() <= maxBounds.Y() {
		return colliders.NewCollision()
	}
	return nil
}
