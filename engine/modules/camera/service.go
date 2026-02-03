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

type Service interface {
	Component() ecs.ComponentsArray[Component]
	Priority() ecs.ComponentsArray[PriorityComponent]

	Mobile() ecs.ComponentsArray[MobileCameraComponent]
	Limits() ecs.ComponentsArray[CameraLimitsComponent]
	Viewport() ecs.ComponentsArray[ViewportComponent]
	NormalizedViewport() ecs.ComponentsArray[NormalizedViewportComponent]

	Ortho() ecs.ComponentsArray[OrthoComponent]
	OrthoResolution() ecs.ComponentsArray[OrthoResolutionComponent]
	Perspective() ecs.ComponentsArray[PerspectiveComponent]
	DynamicPerspective() ecs.ComponentsArray[DynamicPerspectiveComponent]

	// returns cameras from smallest to biggest
	OrderedCameras() []ecs.EntityID

	GetViewport(camera ecs.EntityID) (x, y, w, h int32)
	Mat4(caemra ecs.EntityID) mgl32.Mat4
	ShootRay(camera ecs.EntityID, mousePos window.MousePos) collider.Ray
}
