package transform

import (
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
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

type ParentPivotPointComponent PivotPointComponent
type ParentComponent struct {
	Parent       ecs.EntityID
	RelativeMask uint8
}

const (
	RelativePos uint8 = 1 << iota
	RelativeRotation
	RelativeSize
)

func NewPos(x, y, z float32) PosComponent               { return PosComponent{mgl32.Vec3{x, y, z}} }
func NewRotation(rotation mgl32.Quat) RotationComponent { return RotationComponent{rotation} }
func NewSize(x, y, z float32) SizeComponent             { return SizeComponent{mgl32.Vec3{x, y, z}} }
func NewPivotPoint(x, y, z float32) PivotPointComponent {
	return PivotPointComponent{mgl32.Vec3{x, y, z}}
}
func NewParentPivotPoint(x, y, z float32) ParentPivotPointComponent {
	return ParentPivotPointComponent{mgl32.Vec3{x, y, z}}
}
func NewParent(p ecs.EntityID, mask uint8) ParentComponent { return ParentComponent{p, mask} }
