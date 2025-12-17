package camera

import (
	"engine/modules/collider"
	"engine/services/ecs"
	"engine/services/media/window"
	"errors"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	ErrNotCamera error = errors.New("this isn't a camera")
)

type CameraTool interface {
	Camera() Interface
}

type Interface interface {
	GetObject(ecs.EntityID) (Object, error)
}

//

type Object interface {
	Viewport() (x, y, w, h int32)
	Mat4() mgl32.Mat4
	ShootRay(mousePos window.MousePos) collider.Ray
}
