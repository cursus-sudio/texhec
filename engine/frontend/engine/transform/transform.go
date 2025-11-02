package transform

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	Pos      mgl32.Vec3
	Rotation mgl32.Quat
	Size     mgl32.Vec3
}

func NewTransform() Transform {
	return Transform{
		Pos:      mgl32.Vec3{0, 0, 0},
		Rotation: mgl32.QuatIdent(),
		Size:     mgl32.Vec3{0, 0, 0},
	}
}

func (t Transform) Ptr() *Transform { return &t }
func (t *Transform) Val() Transform { return *t }

func (t1 Transform) Merge(t2 Transform) Transform {
	return Transform{
		Pos:      t1.Pos.Add(t2.Pos),
		Rotation: t1.Rotation.Mul(t2.Rotation),
		Size:     mgl32.Vec3{t1.Size[0] * t2.Size[0], t1.Size[1] * t2.Size[1], t1.Size[2] * t2.Size[2]},
	}
}

func (t *Transform) Mat4() mgl32.Mat4 {
	translation := mgl32.Translate3D(t.Pos.X(), t.Pos.Y(), t.Pos.Z())
	rotation := t.Rotation.Mat4()
	scale := mgl32.Scale3D(t.Size.X()/2, t.Size.Y()/2, t.Size.Z()/2)
	return translation.Mul4(rotation).Mul4(scale)
}

func (t *Transform) SetPos(pos mgl32.Vec3) *Transform {
	t.Pos = pos
	return t
}

func (t *Transform) SetRotation(rotation mgl32.Quat) *Transform {
	t.Rotation = rotation
	return t
}

func (t *Transform) SetSize(size mgl32.Vec3) *Transform {
	t.Size = size
	return t
}
