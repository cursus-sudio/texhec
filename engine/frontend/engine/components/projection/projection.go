package projection

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/groups"

	"github.com/go-gl/mathgl/mgl32"
)

func RayDirection(
	projectionMatrix mgl32.Mat4,
	viewMatrix mgl32.Mat4,
	mousePos mgl32.Vec2,
) (NearWorld mgl32.Vec4, Direction mgl32.Vec3, MaxDistance float32) {
	invViewProjection := projectionMatrix.Mul4(viewMatrix).Inv()

	nearClip := mgl32.Vec4{mousePos.X(), mousePos.Y(), -1, 1}
	farClip := mgl32.Vec4{mousePos.X(), mousePos.Y(), 1, 1}

	nearWorld := invViewProjection.Mul4x1(nearClip)
	farWorld := invViewProjection.Mul4x1(farClip)

	nearWorld = nearWorld.Mul(1 / nearWorld[3])
	farWorld = farWorld.Mul(1 / farWorld[3])

	notNormalizedDirection := farWorld.Vec3().Sub(nearWorld.Vec3())
	direction := notNormalizedDirection.Normalize()
	dist := notNormalizedDirection.Len()

	return nearWorld, direction, dist
}

func ShootRay(
	projectionMatrix mgl32.Mat4,
	viewMatrix mgl32.Mat4,
	mousePos mgl32.Vec2,
	defaultRayOrigin *mgl32.Vec3,
) collider.Ray {
	nearWorld, direction, maxDistance := RayDirection(projectionMatrix, viewMatrix, mousePos)
	var rayOrigin mgl32.Vec3
	if defaultRayOrigin != nil {
		rayOrigin = *defaultRayOrigin
	} else {
		rayOrigin = nearWorld.Vec3()
	}
	return collider.NewRay(rayOrigin, direction, maxDistance, groups.DefaultGroups())
}
