package shapes

import (
	"frontend/services/colliders"
	"frontend/services/graphics/camera"

	"github.com/go-gl/mathgl/mgl32"
)

func (s *collidersService) rect2DRayHandler(rect Rect2D, ray Ray) colliders.Collision {
	planeNormal := rect.Rotation.Rotate(camera.Forward)
	denom := ray.Direction.Dot(planeNormal)
	if mgl32.Abs(denom) < mgl32.Epsilon {
		return nil
	}

	t := rect.Pos.Sub(ray.Pos).Dot(planeNormal) / denom
	if t < 0 {
		return nil
	}

	intersectionPoint := ray.Pos.Add(ray.Direction.Mul(t))

	invRotation := rect.Rotation.Inverse()
	localIntersectionPoint := invRotation.Rotate(intersectionPoint.Sub(rect.Pos))

	halfSizeX, halfSizeY := rect.Size[0]/2, rect.Size[1]/2

	if localIntersectionPoint[0] >= -halfSizeX &&
		localIntersectionPoint[0] <= halfSizeX &&
		localIntersectionPoint[1] >= -halfSizeY &&
		localIntersectionPoint[1] <= halfSizeY {
		return colliders.NewCollision()
	}

	return nil
}
