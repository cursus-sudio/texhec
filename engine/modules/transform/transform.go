package transform

import (
	"engine/modules/hierarchy"
	"engine/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type ToolFactory ecs.ToolFactory[World, TransformTool]
type TransformTool interface {
	Transform() Interface
}
type World interface { // these are dependencies of transform package
	ecs.World
	hierarchy.HierarchyTool
}
type Interface interface {
	SetAbsolutePos(ecs.EntityID, AbsolutePosComponent)
	SetAbsoluteRotation(ecs.EntityID, AbsoluteRotationComponent)
	SetAbsoluteSize(ecs.EntityID, AbsoluteSizeComponent)

	AbsolutePos() ecs.ComponentsArray[AbsolutePosComponent]
	AbsoluteRotation() ecs.ComponentsArray[AbsoluteRotationComponent]
	AbsoluteSize() ecs.ComponentsArray[AbsoluteSizeComponent]

	Pos() ecs.ComponentsArray[PosComponent]
	Rotation() ecs.ComponentsArray[RotationComponent]
	Size() ecs.ComponentsArray[SizeComponent]

	MaxSize() ecs.ComponentsArray[MaxSizeComponent]
	MinSize() ecs.ComponentsArray[MinSizeComponent]

	AspectRatio() ecs.ComponentsArray[AspectRatioComponent]
	PivotPoint() ecs.ComponentsArray[PivotPointComponent]

	Parent() ecs.ComponentsArray[ParentComponent]
	ParentPivotPoint() ecs.ComponentsArray[ParentPivotPointComponent]

	Mat4(ecs.EntityID) mgl32.Mat4
	AddDirtySet(ecs.DirtySet)
}

// components

// parent
type ParentFlag uint8

const (
	RelativePos ParentFlag = 1 << iota
	RelativeRotation
	RelativeSizeX
	RelativeSizeY
	RelativeSizeZ
)
const (
	RelativeSizeXY  = RelativeSizeX | RelativeSizeY
	RelativeSizeXZ  = RelativeSizeX | RelativeSizeZ
	RelativeSizeXYZ = RelativeSizeX | RelativeSizeY | RelativeSizeZ
	RelativeSizeYZ  = RelativeSizeY | RelativeSizeZ
)

// aspect ratio
type PrimaryAxis uint8

const (
	_ PrimaryAxis = iota
	PrimaryAxisX
	PrimaryAxisY
	PrimaryAxisZ
)

type PosComponent struct{ Pos mgl32.Vec3 }
type RotationComponent struct{ Rotation mgl32.Quat }
type SizeComponent struct{ Size mgl32.Vec3 }

type AbsolutePosComponent struct{ Pos mgl32.Vec3 }
type AbsoluteRotationComponent struct{ Rotation mgl32.Quat }
type AbsoluteSizeComponent struct{ Size mgl32.Vec3 }

type MinSizeComponent SizeComponent // refers to absolute size. 0 means ignore axis
type MaxSizeComponent SizeComponent // refers to absolute size. 0 means ignore axis

type AspectRatioComponent struct {
	// 0 means ignore axis
	AspectRatio mgl32.Vec3
	PrimaryAxis PrimaryAxis
}

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
func NewMinSize(x, y, z float32) MinSizeComponent       { return MinSizeComponent{mgl32.Vec3{x, y, z}} }
func NewMaxSize(x, y, z float32) MaxSizeComponent       { return MaxSizeComponent{mgl32.Vec3{x, y, z}} }
func NewAspectRatio(x, y, z float32, primaryAxis PrimaryAxis) AspectRatioComponent {
	return AspectRatioComponent{mgl32.Vec3{x, y, z}, primaryAxis}
}
func NewPivotPoint(x, y, z float32) PivotPointComponent {
	return PivotPointComponent{mgl32.Vec3{x, y, z}}
}
func NewParent(mask ParentFlag) ParentComponent { return ParentComponent{mask} }
func NewParentPivotPoint(x, y, z float32) ParentPivotPointComponent {
	return ParentPivotPointComponent{mgl32.Vec3{x, y, z}}
}

// blend

func (c1 PosComponent) Lerp(c2 PosComponent, mix32 float32) PosComponent {
	return PosComponent{c1.Pos.Mul(1 - mix32).Add(c2.Pos.Mul(mix32))}
}

func (c1 RotationComponent) Lerp(c2 RotationComponent, mix32 float32) RotationComponent {
	return RotationComponent{mgl32.QuatSlerp(c1.Rotation, c2.Rotation, mix32)}
}

func (c1 SizeComponent) Lerp(c2 SizeComponent, mix32 float32) SizeComponent {
	return SizeComponent{c1.Size.Mul(1 - mix32).Add(c2.Size.Mul(mix32))}
}

func (c1 PivotPointComponent) Lerp(c2 PivotPointComponent, mix32 float32) PivotPointComponent {
	return PivotPointComponent{c1.Point.Mul(1 - mix32).Add(c2.Point.Mul(mix32))}
}

func (c1 ParentPivotPointComponent) Lerp(c2 ParentPivotPointComponent, mix32 float32) ParentPivotPointComponent {
	return ParentPivotPointComponent{c1.Point.Mul(1 - mix32).Add(c2.Point.Mul(mix32))}
}
