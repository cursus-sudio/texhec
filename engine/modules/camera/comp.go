package camera

import (
	"reflect"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraComponent struct {
	Projection reflect.Type
}

func NewCamera[Projection any]() CameraComponent {
	return CameraComponent{reflect.TypeFor[Projection]()}
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

//

type ViewportComponent struct{ X, Y, W, H int32 }
type NormalizedViewportComponent struct{ X, Y, W, H float32 }

func NewViewportComponent(x, y, w, h int32) ViewportComponent {
	return ViewportComponent{x, y, w, h}
}
func NewNormalizedViewportComponent(x, y, w, h float32) NormalizedViewportComponent {
	return NormalizedViewportComponent{x, y, w, h}
}

func (c ViewportComponent) Viewport() (x, y, w, h int32) { return c.X, c.Y, c.W, c.H }
func (c NormalizedViewportComponent) Viewport(fullW, fullH int32) (rx, ry, rw, rh int32) { // r is from result
	x := float32(fullW) * c.X
	y := float32(fullH) * c.Y
	w := float32(fullW) * c.W
	h := float32(fullH) * c.H
	return int32(x), int32(y), int32(w), int32(h)
}

//
