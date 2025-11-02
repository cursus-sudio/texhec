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

//

type OrthoComponent struct {
	Width, Height float32
	Near, Far     float32
	Zoom          float32
}

func NewOrtho(w, h, near, far float32, zoom float32) OrthoComponent {
	return OrthoComponent{
		Width:  w / zoom,
		Height: h / zoom,
		Near:   min(near, far),
		Far:    max(near, far),
		Zoom:   zoom,
	}
}

//

type DynamicOrthoComponent struct {
	Near, Far float32
	Zoom      float32
}

func NewDynamicOrtho(near, far float32, zoom float32) DynamicOrthoComponent {
	return DynamicOrthoComponent{
		Near: near,
		Far:  far,
		Zoom: zoom,
	}
}

//

//

type CameraUp mgl32.Vec3
type CameraForward mgl32.Vec3

type Perspective struct {
	FovY        float32
	AspectRatio float32
	Near, Far   float32
}

func NewPerspective(fovY float32, aspectRatio float32, near, far float32) Perspective {
	return Perspective{FovY: fovY, AspectRatio: aspectRatio, Near: near, Far: far}
}

type DynamicPerspective struct {
	FovY      float32
	Near, Far float32
}

func NewDynamicPerspective(fovY float32, near, far float32) DynamicPerspective {
	return DynamicPerspective{
		FovY: fovY,
		Near: near,
		Far:  far,
	}
}
