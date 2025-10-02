package cameras

import (
	"frontend/engine/components/collider"
	"frontend/engine/components/groups"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera interface {
	Mat4() mgl32.Mat4
	ShootRay(mousePos mgl32.Vec2) collider.Ray
}

type camera struct {
	mat4     func() mgl32.Mat4
	shootRay func(mgl32.Vec2) collider.Ray
	groups   groups.Groups
}

func NewCamera(
	mat4 func() mgl32.Mat4,
	shootRay func(mgl32.Vec2) collider.Ray,
	groups groups.Groups,
) Camera {
	return &camera{mat4, shootRay, groups}
}

func (c *camera) Mat4() mgl32.Mat4 { return c.mat4() }
func (c *camera) ShootRay(mousePos mgl32.Vec2) collider.Ray {
	ray := c.shootRay(mousePos)
	ray.Groups = c.groups
	return ray
}
