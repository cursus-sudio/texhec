package transformtool

import (
	"engine/modules/hierarchy"
	"engine/modules/transform"
	"engine/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type object struct {
	transaction

	parent ecs.EntityComponent[hierarchy.ParentComponent]

	pos         ecs.EntityComponent[transform.PosComponent]
	absolutePos ecs.EntityComponent[transform.PosComponent]

	rot         ecs.EntityComponent[transform.RotationComponent]
	absoluteRot ecs.EntityComponent[transform.RotationComponent]

	size         ecs.EntityComponent[transform.SizeComponent]
	absoluteSize ecs.EntityComponent[transform.SizeComponent]

	maxSize ecs.EntityComponent[transform.MaxSizeComponent]
	minSize ecs.EntityComponent[transform.MinSizeComponent]

	pivotPoint       ecs.EntityComponent[transform.PivotPointComponent]
	parentMask       ecs.EntityComponent[transform.ParentComponent]
	parentPivotPoint ecs.EntityComponent[transform.ParentPivotPointComponent]
	entity           ecs.EntityID
}

func newEntityTransform(
	t transaction,
	entity ecs.EntityID,
) transform.Object {
	s := object{
		transaction: t,

		parent: t.parentTransaction.GetEntityComponent(entity),

		pos:  t.posTransaction.GetEntityComponent(entity),
		rot:  t.rotationTransaction.GetEntityComponent(entity),
		size: t.sizeTransaction.GetEntityComponent(entity),

		maxSize: t.maxSizeTransaction.GetEntityComponent(entity),
		minSize: t.minSizeTransaction.GetEntityComponent(entity),

		pivotPoint:       t.pivotPointTransaction.GetEntityComponent(entity),
		parentMask:       t.parentMaskTransaction.GetEntityComponent(entity),
		parentPivotPoint: t.parentPivotPointTransaction.GetEntityComponent(entity),
		entity:           entity,
	}
	s.Init()
	return s
}

func (t object) Pos() ecs.EntityComponent[transform.PosComponent] { return t.pos }
func (t object) AbsolutePos() ecs.EntityComponent[transform.PosComponent] {
	return t.absolutePos
}

func (t object) Rotation() ecs.EntityComponent[transform.RotationComponent] { return t.rot }
func (t object) AbsoluteRotation() ecs.EntityComponent[transform.RotationComponent] {
	return t.absoluteRot
}

func (t object) Size() ecs.EntityComponent[transform.SizeComponent] { return t.size }
func (t object) AbsoluteSize() ecs.EntityComponent[transform.SizeComponent] {
	return t.absoluteSize
}

func (t object) MaxSize() ecs.EntityComponent[transform.MaxSizeComponent] { return t.maxSize }
func (t object) MinSize() ecs.EntityComponent[transform.MinSizeComponent] { return t.minSize }

func (t object) PivotPoint() ecs.EntityComponent[transform.PivotPointComponent] {
	return t.pivotPoint
}

func (t object) Parent() ecs.EntityComponent[transform.ParentComponent] { return t.parentMask }
func (t object) ParentPivotPoint() ecs.EntityComponent[transform.ParentPivotPointComponent] {
	return t.parentPivotPoint
}

func (t object) Mat4() mgl32.Mat4 {
	pos, err := t.absolutePos.Get()
	if err != nil {
		pos = transform.NewPos(0, 0, 0)
	}
	rot, err := t.absoluteRot.Get()
	if err != nil {
		rot = transform.NewRotation(mgl32.QuatIdent())
	}
	size, err := t.absoluteSize.Get()
	if err != nil {
		size = transform.NewSize(1, 1, 1)
	}

	translation := mgl32.Translate3D(pos.Pos.X(), pos.Pos.Y(), pos.Pos.Z())
	rotation := rot.Rotation.Mat4()
	scale := mgl32.Scale3D(size.Size.X()/2, size.Size.Y()/2, size.Size.Z()/2)
	return translation.Mul4(rotation).Mul4(scale)
}
