package camera

import (
	"github.com/go-gl/mathgl/mgl32"
)

type OrthoComponent struct {
	Near, Far float32
	Zoom      float32
}

func NewOrtho(near, far float32) OrthoComponent {
	return OrthoComponent{
		Near: near,
		Far:  far,
		Zoom: 1,
	}
}

func (c OrthoComponent) GetMatrix(w, h int32) mgl32.Mat4 {
	fW, fH := float32(w), float32(h)

	return mgl32.Ortho(
		-fW/c.Zoom/2, fW/c.Zoom/2,
		-fH/c.Zoom/2, fH/c.Zoom/2,
		c.Near, c.Far,
	)
}

//

type OrthoResolutionComponent struct {
	W, H int32
}

func NewOrthoResolution(w, h int32) OrthoResolutionComponent { return OrthoResolutionComponent{w, h} }
func GetViewportOrthoResolution(x, y, w, h int32) OrthoResolutionComponent {
	return OrthoResolutionComponent{w - x, h - y}
}
func (c *OrthoResolutionComponent) Elem() (w, h int32) { return c.W, c.H }

//

type CameraUp mgl32.Vec3
type CameraForward mgl32.Vec3

type PerspectiveComponent struct {
	FovY        float32
	AspectRatio float32
	Near, Far   float32
}

func NewPerspective(fovY float32, aspectRatio float32, near, far float32) PerspectiveComponent {
	return PerspectiveComponent{FovY: fovY, AspectRatio: aspectRatio, Near: near, Far: far}
}

type DynamicPerspectiveComponent struct {
	FovY      float32
	Near, Far float32
}

func NewDynamicPerspective(fovY float32, near, far float32) DynamicPerspectiveComponent {
	return DynamicPerspectiveComponent{
		FovY: fovY,
		Near: near,
		Far:  far,
	}
}
