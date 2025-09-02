package projection

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/transform"
	"frontend/services/graphics/camera"

	"github.com/go-gl/mathgl/mgl32"
)

type Perspective struct {
	FovY        float32
	AspectRatio float32
	Near, Far   float32
}

func NewPerspective(fovY float32, aspectRatio float32, near, far float32) Perspective {
	return Perspective{FovY: fovY, AspectRatio: aspectRatio, Near: near, Far: far}
}

func (p Perspective) Mat4() mgl32.Mat4 {
	return mgl32.Perspective(p.FovY, p.AspectRatio, p.Near, p.Far)
}

func (p Perspective) ViewMat4(cameraTransform transform.Transform) mgl32.Mat4 {
	return mgl32.LookAtV(
		cameraTransform.Pos,
		cameraTransform.Pos.Add(cameraTransform.Rotation.Rotate(camera.Forward)),
		camera.Up,
	)
}

func (p Perspective) ShootRay(
	cameraTransform transform.Transform,
	mousePos mgl32.Vec2,
) collider.Ray {
	return ShootRay(
		p.Mat4(),
		p.ViewMat4(cameraTransform),
		mousePos,
		&cameraTransform.Pos,
	)
}
