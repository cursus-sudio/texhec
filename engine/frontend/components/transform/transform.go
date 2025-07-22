package transform

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	Pos      Pos
	Rotation mgl32.Quat
	Size     Size
}

func NewTransform() Transform {
	return Transform{
		Pos:      Pos{X: 0, Y: 0, Z: 0},
		Rotation: mgl32.QuatIdent(),
		Size:     Size{X: 0, Y: 0, Z: 0},
	}
}

func (t Transform) SetPos(pos Pos) Transform {
	return Transform{Pos: pos, Rotation: t.Rotation, Size: t.Size}
}

func (t Transform) SetRotation(rotation mgl32.Quat) Transform {
	return Transform{Pos: t.Pos, Rotation: rotation, Size: t.Size}
}

func (t Transform) SetSize(size Size) Transform {
	return Transform{Pos: t.Pos, Rotation: t.Rotation, Size: size}
}
