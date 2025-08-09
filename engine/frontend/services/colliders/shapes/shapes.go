package shapes

import (
	"frontend/engine/components/transform"
	"frontend/services/colliders"
	"frontend/services/graphics/camera"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

// ray

type ray struct {
	Pos       mgl32.Vec3
	Direction mgl32.Vec3
}

type Ray struct{ ray }

func (s ray) Apply(t transform.Transform) colliders.Shape {
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

// func rayRayHandler(s1 Ray, s2 Ray) colliders.Collision {
// 	return nil
// }
