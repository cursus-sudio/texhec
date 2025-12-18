package camera

import (
	"engine/modules/collider"
	"engine/modules/transform"
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

type World interface {
	ecs.World
	transform.TransformTool
}

type Interface interface {
	GetObject(ecs.EntityID) (Object, error)
	Component() ecs.ComponentsArray[Component]

	Mobile() ecs.ComponentsArray[MobileCameraComponent]
	Limits() ecs.ComponentsArray[CameraLimitsComponent]
	Viewport() ecs.ComponentsArray[ViewportComponent]
	NormalizedViewport() ecs.ComponentsArray[NormalizedViewportComponent]

	Ortho() ecs.ComponentsArray[OrthoComponent]
	OrthoResolution() ecs.ComponentsArray[OrthoResolutionComponent]
	Perspective() ecs.ComponentsArray[PerspectiveComponent]
	DynamicPerspective() ecs.ComponentsArray[DynamicPerspectiveComponent]
}

//

type Object interface {
	Viewport() (x, y, w, h int32)
	Mat4() mgl32.Mat4
	ShootRay(mousePos window.MousePos) collider.Ray
}
