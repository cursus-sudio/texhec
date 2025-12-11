package transformtool

import (
	"engine/modules/hierarchy"
	"engine/modules/transform"
	"engine/services/ecs"
)

type transaction struct {
	tool

	parentTransaction           ecs.ComponentsArrayTransaction[hierarchy.ParentComponent]
	posTransaction              ecs.ComponentsArrayTransaction[transform.PosComponent]
	rotationTransaction         ecs.ComponentsArrayTransaction[transform.RotationComponent]
	sizeTransaction             ecs.ComponentsArrayTransaction[transform.SizeComponent]
	maxSizeTransaction          ecs.ComponentsArrayTransaction[transform.MaxSizeComponent]
	minSizeTransaction          ecs.ComponentsArrayTransaction[transform.MinSizeComponent]
	aspectRatioTransaction      ecs.ComponentsArrayTransaction[transform.AspectRatioComponent]
	pivotPointTransaction       ecs.ComponentsArrayTransaction[transform.PivotPointComponent]
	parentMaskTransaction       ecs.ComponentsArrayTransaction[transform.ParentComponent]
	parentPivotPointTransaction ecs.ComponentsArrayTransaction[transform.ParentPivotPointComponent]
}

func newTransformTransaction(
	tool tool,
) transform.Transaction {
	return transaction{
		tool,
		tool.parentArray.Transaction(),
		tool.posArray.Transaction(),
		tool.rotationArray.Transaction(),
		tool.sizeArray.Transaction(),
		tool.maxSizeArray.Transaction(),
		tool.minSizeArray.Transaction(),
		tool.aspectRatioArray.Transaction(),
		tool.pivotPointArray.Transaction(),
		tool.parentMaskArray.Transaction(),
		tool.parentPivotPointArray.Transaction(),
	}
}

func (t transaction) GetObject(entity ecs.EntityID) transform.Object {
	return newEntityTransform(t, entity)
}

func (t transaction) Transactions() []ecs.AnyComponentsArrayTransaction {
	return []ecs.AnyComponentsArrayTransaction{
		t.posTransaction,
		t.rotationTransaction,
		t.sizeTransaction,
		t.maxSizeTransaction,
		t.minSizeTransaction,
		t.aspectRatioTransaction,
		t.pivotPointTransaction,
		t.parentMaskTransaction,
		t.parentPivotPointTransaction,
	}
}

func (t transaction) Flush() error {
	return ecs.FlushMany(t.Transactions()...)
}
