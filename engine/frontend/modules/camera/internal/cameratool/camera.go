package cameratool

import (
	"frontend/modules/camera"
	"frontend/modules/collider"
	"frontend/modules/groups"
	"frontend/services/media/window"

	"github.com/go-gl/mathgl/mgl32"
)

type cameraService struct {
	mat4     func() mgl32.Mat4
	viewport func() (x, y, w, h int32)
	shootRay func(window.MousePos) collider.Ray
	groups   groups.GroupsComponent
}

func NewCameraService(
	mat4 func() mgl32.Mat4,
	viewport func() (x, y, w, h int32),
	shootRay func(mousePos window.MousePos) collider.Ray,
	groups groups.GroupsComponent,
) camera.CameraService {
	return &cameraService{mat4, viewport, shootRay, groups}
}

func (c *cameraService) Mat4() mgl32.Mat4             { return c.mat4() }
func (c *cameraService) Viewport() (x, y, w, h int32) { return c.viewport() }
func (c *cameraService) ShootRay(mousePos window.MousePos) collider.Ray {
	ray := c.shootRay(mousePos)
	ray.Groups = c.groups
	return ray
}
