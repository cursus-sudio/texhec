package transform

import "github.com/go-gl/mathgl/mgl32"

// lock refers to object center.
// default center is (.5, .5, .5).
// each axis value should be between 0 and 1.
//
// example: to align to left use (0, .5, .5)
type PosLock struct {
	Lock mgl32.Vec3
}

func NewPosLock(lock mgl32.Vec3) PosLock {
	return PosLock{lock}
}
