package shapes

import (
	"frontend/engine/components/transform"
	"frontend/services/colliders"

	"github.com/go-gl/mathgl/mgl32"
)

// ray

type ray struct {
	transform.Transform
}

type Ray struct{ ray }

func (s ray) Apply(t transform.Transform) colliders.Shape {
	return Ray{ray{s.Transform.Merge(t)}}
}

func (s ray) Position() mgl32.Vec3 { return s.Pos }

func NewRay(pos mgl32.Vec3, rotation mgl32.Quat) Ray {
	return Ray{ray{
		transform.NewTransform().
			SetPos(pos).
			SetRotation(rotation),
	}}
}

// func rayRayHandler(s1 Ray, s2 Ray) colliders.Collision {
// 	return nil
// }
