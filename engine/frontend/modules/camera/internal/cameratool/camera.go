package cameratool

import (
	"frontend/modules/camera"
	"frontend/modules/collider"
	"frontend/modules/groups"

	"github.com/go-gl/mathgl/mgl32"
)

type cameraService struct {
	mat4     func() mgl32.Mat4
	shootRay func(mgl32.Vec2) collider.Ray
	groups   groups.GroupsComponent
}

func NewCameraService(
	mat4 func() mgl32.Mat4,
	shootRay func(mgl32.Vec2) collider.Ray,
	groups groups.GroupsComponent,
) camera.CameraService {
	return &cameraService{mat4, shootRay, groups}
}

func (c *cameraService) Mat4() mgl32.Mat4 { return c.mat4() }
func (c *cameraService) ShootRay(mousePos mgl32.Vec2) collider.Ray {
	ray := c.shootRay(mousePos)
	ray.Groups = c.groups
	return ray
}
