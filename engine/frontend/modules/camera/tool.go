package camera

import (
	"errors"
	"frontend/modules/collider"
	"frontend/services/media/window"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	ErrNotCamera error = errors.New("this isn't a camera")
)

type Tool interface {
	GetObject(ecs.EntityID) (Object, error)
}

//

type Object interface {
	Viewport() (x, y, w, h int32)
	Mat4() mgl32.Mat4
	ShootRay(mousePos window.MousePos) collider.Ray
}
