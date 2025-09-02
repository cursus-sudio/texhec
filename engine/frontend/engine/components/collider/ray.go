package collider

import (
	"frontend/engine/components/transform"
	"frontend/services/graphics/camera"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type ray struct {
	Pos       mgl32.Vec3
	Direction mgl32.Vec3
}

type Ray struct{ ray }

func (s ray) Apply(t transform.Transform) Ray {
	return Ray{ray{Pos: s.Pos.Add(t.Pos), Direction: t.Rotation.Rotate(s.Direction)}}
}

func (s ray) Position() mgl32.Vec3 { return s.Pos }

func (s ray) Rotation() mgl32.Quat {
	axis := camera.Forward.Cross(s.Direction).Normalize()
	angle := float32(math.Acos(float64(camera.Forward.Dot(s.Direction))))

	return mgl32.QuatRotate(angle, axis)
}

func NewRay(pos mgl32.Vec3, direction mgl32.Vec3) Ray {
	return Ray{ray{
		Pos:       pos,
		Direction: direction.Normalize(),
	}}
}

type RayHit struct {
	Point    mgl32.Vec3
	Normal   mgl32.Vec3
	Distance float32
}

func NewRayHit(point mgl32.Vec3, normal mgl32.Vec3, dist float32) RayHit {
	return RayHit{point, normal, dist}
}
