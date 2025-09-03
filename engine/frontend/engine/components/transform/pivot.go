package transform

import "github.com/go-gl/mathgl/mgl32"

// pivot refers to object center.
// default center is (.5, .5, .5).
// each axis value should be between 0 and 1.
//
// example: to align to left use (0, .5, .5)
type PivotPoint struct {
	Lock mgl32.Vec3
}

func NewPivotPoint(lock mgl32.Vec3) PivotPoint {
	return PivotPoint{lock}
}
