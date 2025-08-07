package transform

import (
	"github.com/go-gl/mathgl/mgl32"
)

var Up mgl32.Vec3 = mgl32.Vec3{0, 1, 0}
var Foward mgl32.Vec3 = mgl32.Vec3{0, 0, -1}

type transform struct {
	Pos      mgl32.Vec3
	Rotation mgl32.Quat
	Size     mgl32.Vec3
}

type Transform struct{ *transform }

func NewTransform() Transform {
	return Transform{&transform{
		Pos:      mgl32.Vec3{0, 0, 0},
		Rotation: mgl32.QuatIdent(),
		Size:     mgl32.Vec3{0, 0, 0},
	}}
}

func (t1 Transform) Merge(t2 Transform) Transform {
	return Transform{&transform{
		Pos:      t1.Pos.Add(t2.Pos),
		Rotation: t1.Rotation.Mul(t2.Rotation),
		Size:     mgl32.Vec3{t1.Size[0] * t2.Size[0], t1.Size[1] * t2.Size[1], t1.Size[2] * t2.Size[2]},
	}}
}

func (t *Transform) Mat4() mgl32.Mat4 {
	translation := mgl32.Translate3D(t.Pos.X(), t.Pos.Y(), t.Pos.Z())
	rotation := t.Rotation.Mat4()
	scale := mgl32.Scale3D(t.Size.X()/2, t.Size.Y()/2, t.Size.Z()/2)
	return translation.Mul4(rotation).Mul4(scale)
}

func (t Transform) SetPos(pos mgl32.Vec3) Transform {
	return Transform{&transform{Pos: pos, Rotation: t.Rotation, Size: t.Size}}
}

func (t Transform) SetRotation(rotation mgl32.Quat) Transform {
	return Transform{&transform{Pos: t.Pos, Rotation: rotation, Size: t.Size}}
}

func (t Transform) SetSize(size mgl32.Vec3) Transform {
	return Transform{&transform{Pos: t.Pos, Rotation: t.Rotation, Size: size}}
}

type aabb struct {
	Min, Max mgl32.Vec3
}

type AABB struct{ *aabb }

func NewAABB(min, max mgl32.Vec3) AABB {
	return AABB{&aabb{Min: min, Max: max}}
}

func (t Transform) ToAABB() AABB {
	halfSize := t.Size.Mul(0.5)

	corners := [8]mgl32.Vec3{
		{-1, -1, -1}, {1, -1, -1}, {-1, 1, -1}, {1, 1, -1},
		{-1, -1, 1}, {1, -1, 1}, {-1, 1, 1}, {1, 1, 1},
	}

	var minCorner, maxCorner mgl32.Vec3

	for i, corner := range corners {
		transformedCorner := t.Rotation.
			Rotate(mgl32.Vec3{corner[0] * halfSize[0], corner[1] * halfSize[1], corner[2] * halfSize[2]}).
			Add(t.Pos)

		if i == 0 {
			minCorner, maxCorner = transformedCorner, transformedCorner
			continue
		}
		minCorner = mgl32.Vec3{
			min(minCorner[0], transformedCorner[0]),
			min(minCorner[1], transformedCorner[1]),
			min(minCorner[2], transformedCorner[2]),
		}
		maxCorner = mgl32.Vec3{
			max(maxCorner[0], transformedCorner[0]),
			max(maxCorner[1], transformedCorner[1]),
			max(maxCorner[2], transformedCorner[2]),
		}
	}

	return NewAABB(minCorner, maxCorner)
}
