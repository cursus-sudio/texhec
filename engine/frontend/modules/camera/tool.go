package camera

import (
	"errors"
	"frontend/modules/collider"
	"frontend/services/media/window"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraService interface {
	Mat4() mgl32.Mat4
	Viewport() (x, y, w, h int32)
	ShootRay(mousePos window.MousePos) collider.Ray
}

//

var (
	ErrNotCamera error = errors.New("this isn't a camera")
)

type CameraTool interface {
	Get(ecs.EntityID) (CameraService, error)
}

//
