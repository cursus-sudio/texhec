package colliders

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Intersection interface {
	PointOnA() mgl32.Vec3
	PointOnB() mgl32.Vec3
	Normal() mgl32.Vec3
	Depth() float32
	Reverse() Intersection
}

type intersection struct {
	pointOnA,
	pointOnB mgl32.Vec3
	normal mgl32.Vec3
	depth  float32
}

func (intersection *intersection) PointOnA() mgl32.Vec3 { return intersection.pointOnA }
func (intersection *intersection) PointOnB() mgl32.Vec3 { return intersection.pointOnB }
func (intersection *intersection) Normal() mgl32.Vec3   { return intersection.normal }
func (intersection *intersection) Depth() float32       { return intersection.depth }
func (intersection *intersection) Reverse() Intersection {
	return NewIntersection(
		intersection.pointOnB,
		intersection.pointOnA,
		mgl32.Vec3{-intersection.normal.X(), -intersection.normal.Y(), -intersection.normal.Z()},
		intersection.depth,
	)
}

func NewIntersection(pointOnA, pointOnB mgl32.Vec3, normal mgl32.Vec3, depth float32) Intersection {
	return &intersection{
		pointOnA: pointOnA,
		pointOnB: pointOnB,
		normal:   normal,
		depth:    depth,
	}
}
