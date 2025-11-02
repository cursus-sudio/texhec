package camera

import (
	"errors"
	"frontend/modules/collider"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraService interface {
	Mat4() mgl32.Mat4
	ShootRay(mousePos mgl32.Vec2) collider.Ray
}

//

var (
	ErrNotCamera error = errors.New("this isn't a camera")
)

type CameraTool interface {
	Get(ecs.EntityID) (CameraService, error)
}

//
