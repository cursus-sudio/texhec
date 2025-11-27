package transform

import (
	"github.com/go-gl/mathgl/mgl32"
)

// components

type ParentFlag uint8

const (
	RelativePos ParentFlag = 1 << iota
	RelativeRotation
	RelativeSize
)

type PosComponent struct{ Pos mgl32.Vec3 }
type RotationComponent struct{ Rotation mgl32.Quat }
type SizeComponent struct{ Size mgl32.Vec3 }

// pivot refers to object center.
// default center is (.5, .5, .5).
// each axis value should be between 0 and 1.
//
// example: to align to left use (0, .5, .5)
type PivotPointComponent struct{ Point mgl32.Vec3 }

type ParentComponent struct{ RelativeMask ParentFlag }
type ParentPivotPointComponent PivotPointComponent

// ctors

func NewPos(x, y, z float32) PosComponent               { return PosComponent{mgl32.Vec3{x, y, z}} }
func NewRotation(rotation mgl32.Quat) RotationComponent { return RotationComponent{rotation} }
func NewSize(x, y, z float32) SizeComponent             { return SizeComponent{mgl32.Vec3{x, y, z}} }
func NewPivotPoint(x, y, z float32) PivotPointComponent {
	return PivotPointComponent{mgl32.Vec3{x, y, z}}
}
func NewParent(mask ParentFlag) ParentComponent { return ParentComponent{mask} }
func NewParentPivotPoint(x, y, z float32) ParentPivotPointComponent {
	return ParentPivotPointComponent{mgl32.Vec3{x, y, z}}
}

// blend

func (c1 PosComponent) Blend(c2 PosComponent, mix32 float32) PosComponent {
	return PosComponent{c1.Pos.Mul(1 - mix32).Add(c2.Pos.Mul(mix32))}
}

func (c1 RotationComponent) Blend(c2 RotationComponent, mix32 float32) RotationComponent {
	return RotationComponent{mgl32.QuatSlerp(c1.Rotation, c2.Rotation, mix32)}
}

func (c1 SizeComponent) Blend(c2 SizeComponent, mix32 float32) SizeComponent {
	return SizeComponent{c1.Size.Mul(1 - mix32).Add(c2.Size.Mul(mix32))}
}

func (c1 PivotPointComponent) Blend(c2 PivotPointComponent, mix32 float32) PivotPointComponent {
	return PivotPointComponent{c1.Point.Mul(1 - mix32).Add(c2.Point.Mul(mix32))}
}

func (c1 ParentPivotPointComponent) Blend(c2 ParentPivotPointComponent, mix32 float32) ParentPivotPointComponent {
	return ParentPivotPointComponent{c1.Point.Mul(1 - mix32).Add(c2.Point.Mul(mix32))}
}
