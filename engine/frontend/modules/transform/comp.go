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

func NewPos(pos mgl32.Vec3) PosComponent                         { return PosComponent{pos} }
func NewRotation(rotation mgl32.Quat) RotationComponent          { return RotationComponent{rotation} }
func NewSize(size mgl32.Vec3) SizeComponent                      { return SizeComponent{size} }
func NewPivotPoint(point mgl32.Vec3) PivotPointComponent         { return PivotPointComponent{point} }
func NewParentPivotPoint(p mgl32.Vec3) ParentPivotPointComponent { return ParentPivotPointComponent{p} }
func NewParent(p ecs.EntityID, mask uint8) ParentComponent       { return ParentComponent{p, mask} }
