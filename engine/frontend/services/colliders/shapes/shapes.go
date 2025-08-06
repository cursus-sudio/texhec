package shapes

import (
	"frontend/engine/components/transform"
	"frontend/services/colliders"

	"github.com/go-gl/mathgl/mgl32"
)

// ray

type ray struct {
	pos      mgl32.Vec3
	rotation mgl32.Quat // or Vec3 is also good because ray doesn't have orientation only direction
}

type Ray struct{ ray }

func (s ray) Apply(t transform.Transform) colliders.Shape {
	return Ray{ray{s.pos.Add(t.Pos), s.rotation.Mul(t.Rotation)}}
}

func (s ray) Position() mgl32.Vec3 { return s.pos }

func NewRay(pos mgl32.Vec3, rotation mgl32.Quat) Ray {
	return Ray{ray{pos, rotation}}
}

// func rayRayHandler(s1 Ray, s2 Ray) colliders.Collision {
// 	return nil
// }
