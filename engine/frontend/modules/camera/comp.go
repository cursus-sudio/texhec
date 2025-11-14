package camera

import (
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraComponent struct {
	Projection ecs.ComponentType
}

func NewCamera(projection ecs.ComponentType) CameraComponent {
	return CameraComponent{projection}
}

//

// component specifying that camera can be freely moved on map
type MobileCameraComponent struct{}

func NewMobileCamera() MobileCameraComponent { return MobileCameraComponent{} }

//

type CameraLimitsComponent struct{ Min, Max mgl32.Vec3 }

func NewCameraLimits(min, max mgl32.Vec3) CameraLimitsComponent {
	return CameraLimitsComponent{min, max}
}
