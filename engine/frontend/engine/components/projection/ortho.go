package projection

import (
	"frontend/engine/components/transform"
	"frontend/services/colliders/shapes"

	"github.com/go-gl/mathgl/mgl32"
)

type Ortho struct {
	Width, Height float32
	Near, Far     float32
}

func NewOrtho(w, h, near, far float32) Ortho {
	return Ortho{Width: w, Height: h, Near: near, Far: far}
}

func (ortho Ortho) Mat4() mgl32.Mat4 {
	return mgl32.Ortho(
		-ortho.Width/2, ortho.Width/2,
		-ortho.Height/2, ortho.Height/2,
		ortho.Near, ortho.Far,
	)
}

func (p Ortho) ViewMat4(cameraTransform transform.Transform) mgl32.Mat4 {
	cameraRotation := cameraTransform.Rotation.Inverse()
	cameraPosition := cameraTransform.Rotation.Rotate(cameraTransform.Pos.Mul(-1))
	return cameraRotation.Mat4().Mul4(mgl32.Translate3D(cameraPosition.X(), cameraPosition.Y(), cameraPosition.Z()))
}

func (ortho Ortho) ShootRay(
	cameraTransform transform.Transform,
	mousePos mgl32.Vec2,
) shapes.Ray {
	return ShootRay(
		ortho.Mat4(),
		ortho.ViewMat4(cameraTransform),
		mousePos,
		nil,
	)
}
