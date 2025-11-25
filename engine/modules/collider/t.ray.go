package collider

import (
	"engine/modules/groups"

	"github.com/go-gl/mathgl/mgl32"
)

type Ray struct {
	Pos       mgl32.Vec3
	Direction mgl32.Vec3
	// max length is either 0 symbolizing infinity or a potive number
	MaxDistance float32
	Groups      groups.GroupsComponent // collision mask
}

func NewRay(pos mgl32.Vec3, direction mgl32.Vec3, maxDistance float32, groups groups.GroupsComponent) Ray {
	return Ray{
		Pos:         pos,
		Direction:   direction.Normalize(),
		MaxDistance: max(0, maxDistance),
		Groups:      groups,
	}
}

// TODO think about moving to tool

func (r *Ray) Apply(transform mgl32.Mat4) {
	newPos := transform.Mul4x1(r.Pos.Vec4(1.0)).Vec3()
	r.Pos = newPos

	newDirection := transform.Transpose().Mat3().Mul3x1(r.Direction)
	scaleFactor := newDirection.Len()
	if scaleFactor < mgl32.Epsilon {
		newDirection = r.Direction
		scaleFactor = 1
	}
	r.Direction = newDirection.Normalize()

	newDirection = newDirection.Normalize()

	var newMaxDistance float32
	if r.MaxDistance == 0.0 {
		newMaxDistance = 0.0
	} else {
		newMaxDistance = r.MaxDistance * scaleFactor
	}
	r.MaxDistance = newMaxDistance
}

func (r Ray) HitPoint() mgl32.Vec3 { return r.Pos.Add(r.Direction.Mul(r.MaxDistance)) }

//

type RayHit struct {
	Point    mgl32.Vec3
	Normal   mgl32.Vec3
	Distance float32
}

func NewRayHit(ray Ray, normal mgl32.Vec3) RayHit {
	return RayHit{ray.HitPoint(), normal, ray.MaxDistance}
}
