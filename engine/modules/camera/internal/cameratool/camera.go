package cameratool

import (
	"engine/modules/camera"
	"engine/modules/collider"
	"engine/modules/groups"
	"engine/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
)

type object struct {
	mat4     func() mgl32.Mat4
	viewport func() (x, y, w, h int32)
	shootRay func(window.MousePos) collider.Ray
	groups   groups.GroupsComponent
}

func NewObject(
	mat4 func() mgl32.Mat4,
	viewport func() (x, y, w, h int32),
	shootRay func(mousePos window.MousePos) collider.Ray,
	groups groups.GroupsComponent,
) camera.Object {
	return &object{mat4, viewport, shootRay, groups}
}

func (c *object) Mat4() mgl32.Mat4             { return c.mat4() }
func (c *object) Viewport() (x, y, w, h int32) { return c.viewport() }
func (c *object) ShootRay(mousePos window.MousePos) collider.Ray {
	ray := c.shootRay(mousePos)
	ray.Groups = c.groups
	return ray
}
