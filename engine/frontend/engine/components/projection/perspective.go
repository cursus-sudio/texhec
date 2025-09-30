package projection

import "github.com/go-gl/mathgl/mgl32"

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
