package collider

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Ray struct {
	Pos       mgl32.Vec3
	Direction mgl32.Vec3
	// max length is either 0 symbolizing infinity or a potive number
	MaxDistance float32
}

func NewRay(pos mgl32.Vec3, direction mgl32.Vec3, maxDistance float32) Ray {
	return Ray{
		Pos:         pos,
		Direction:   direction.Normalize(),
		MaxDistance: max(0, maxDistance),
	}
}

func (r Ray) HitPoint() mgl32.Vec3 { return r.Pos.Add(r.Direction.Mul(r.MaxDistance)) }

type RayHit struct {
	Point    mgl32.Vec3
	Normal   mgl32.Vec3
	Distance float32
}

func NewRayHit(ray Ray, normal mgl32.Vec3) RayHit {
	return RayHit{ray.HitPoint(), normal, ray.MaxDistance}
}
