package cameras

import (
	"frontend/engine/components/collider"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera interface {
	Mat4() mgl32.Mat4
	ShootRay(mousePos mgl32.Vec2) collider.Ray
}

type camera struct {
	mat4     func() mgl32.Mat4
	shootRay func(mgl32.Vec2) collider.Ray
}

func NewCamera(
	mat4 func() mgl32.Mat4,
	shootRay func(mgl32.Vec2) collider.Ray,
) Camera {
	return &camera{mat4, shootRay}
}

func (c *camera) Mat4() mgl32.Mat4                          { return c.mat4() }
func (c *camera) ShootRay(mousePos mgl32.Vec2) collider.Ray { return c.shootRay(mousePos) }
