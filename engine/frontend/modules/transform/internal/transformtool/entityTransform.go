package transformtool

import (
	"frontend/modules/transform"
	"shared/services/datastructures"
	"shared/services/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type entityTransform struct {
	transformTransaction

	pos         ecs.EntityComponent[transform.PosComponent]
	absolutePos ecs.EntityComponent[transform.PosComponent]

	rot         ecs.EntityComponent[transform.RotationComponent]
	absoluteRot ecs.EntityComponent[transform.RotationComponent]

	size         ecs.EntityComponent[transform.SizeComponent]
	absoluteSize ecs.EntityComponent[transform.SizeComponent]

	pivotPoint       ecs.EntityComponent[transform.PivotPointComponent]
	parent           ecs.EntityComponent[transform.ParentComponent]
	parentPivotPoint ecs.EntityComponent[transform.ParentPivotPointComponent]
	entity           ecs.EntityID
}

func newEntityTransform(
	t transformTransaction,
	entity ecs.EntityID,
) transform.EntityTransform {
	s := entityTransform{
		transformTransaction: t,

		pos:  t.posArray.GetEntityComponent(entity, t.posTransaction),
		rot:  t.rotationArray.GetEntityComponent(entity, t.rotationTransaction),
		size: t.sizeArray.GetEntityComponent(entity, t.sizeTransaction),

		pivotPoint:       t.pivotPointArray.GetEntityComponent(entity, t.pivotPointTransaction),
		parent:           t.parentArray.GetEntityComponent(entity, t.parentTransaction),
		parentPivotPoint: t.parentPivotPointArray.GetEntityComponent(entity, t.parentPivotPointTransaction),
		entity:           entity,
	}
	s.Init()
	return s
}

func (t entityTransform) Pos() ecs.EntityComponent[transform.PosComponent] { return t.pos }
func (t entityTransform) AbsolutePos() ecs.EntityComponent[transform.PosComponent] {
	return t.absolutePos
}

func (t entityTransform) Rotation() ecs.EntityComponent[transform.RotationComponent] { return t.rot }
func (t entityTransform) AbsoluteRotation() ecs.EntityComponent[transform.RotationComponent] {
	return t.absoluteRot
}

func (t entityTransform) Size() ecs.EntityComponent[transform.SizeComponent] { return t.size }
func (t entityTransform) AbsoluteSize() ecs.EntityComponent[transform.SizeComponent] {
	return t.absoluteSize
}

func (t entityTransform) PivotPoint() ecs.EntityComponent[transform.PivotPointComponent] {
	return t.pivotPoint
}

func (t entityTransform) Parent() ecs.EntityComponent[transform.ParentComponent] { return t.parent }
func (t entityTransform) ParentPivotPoint() ecs.EntityComponent[transform.ParentPivotPointComponent] {
	return t.parentPivotPoint
}

func (t entityTransform) Children() datastructures.SparseSetReader[ecs.EntityID] {
	return t.parentTool.GetChildren(t.entity)
}

func (t entityTransform) FlatChildren() datastructures.SparseSetReader[ecs.EntityID] {
	flatChildren := datastructures.NewSparseSet[ecs.EntityID]()
	leftParents := []ecs.EntityID{t.entity}
	for len(leftParents) != 0 {
		parent := leftParents[0]
		leftParents = leftParents[1:]
		parentTransform := t.GetEntity(parent)
		childrenSet := parentTransform.Children()
		children := childrenSet.GetIndices()
		for _, child := range children {
			flatChildren.Add(child)
		}
		leftParents = append(leftParents, children...)
	}
	return flatChildren
}

func (t entityTransform) Mat4() mgl32.Mat4 {
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
