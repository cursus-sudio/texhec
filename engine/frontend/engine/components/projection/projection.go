package projection

import (
	"frontend/engine/components/transform"
	"frontend/services/colliders/shapes"

	"github.com/go-gl/mathgl/mgl32"
)

type Projection interface {
	Mat4() mgl32.Mat4
	ViewMat4(transform.Transform) mgl32.Mat4
	ShootRay(
		cameraTransform transform.Transform,
		mousePos mgl32.Vec2,
	) shapes.Ray
}

func RayDirection(
	projectionMatrix mgl32.Mat4,
	viewMatrix mgl32.Mat4,
	mousePos mgl32.Vec2,
) (NearWorld mgl32.Vec4, Direction mgl32.Vec3) {
	invViewProjection := projectionMatrix.Mul4(viewMatrix).Inv()

	nearClip := mgl32.Vec4{mousePos.X(), mousePos.Y(), -1, 1}
	farClip := mgl32.Vec4{mousePos.X(), mousePos.Y(), 1, 1}

	nearWorld := invViewProjection.Mul4x1(nearClip)
	farWorld := invViewProjection.Mul4x1(farClip)

	nearWorld = nearWorld.Mul(1 / nearWorld[3])
	farWorld = farWorld.Mul(1 / farWorld[3])

	direction := mgl32.Vec3{
		farWorld[0] - nearWorld[0],
		farWorld[1] - nearWorld[1],
		farWorld[2] - nearWorld[2],
	}.Normalize()

	return nearWorld, direction
}

func ShootRay(
	projectionMatrix mgl32.Mat4,
	viewMatrix mgl32.Mat4,
	mousePos mgl32.Vec2,
	defaultRayOrigin *mgl32.Vec3,
) shapes.Ray {
	nearWorld, direction := RayDirection(projectionMatrix, viewMatrix, mousePos)
	var rayOrigin mgl32.Vec3
	if defaultRayOrigin != nil {
		rayOrigin = *defaultRayOrigin
	} else {
		rayOrigin = nearWorld.Vec3()
	}
	return shapes.NewRay(rayOrigin, direction)
}
