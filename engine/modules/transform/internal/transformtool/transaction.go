package transformtool

import (
	"engine/modules/transform"
	"engine/services/ecs"
)

type transaction struct {
	tool

	posTransaction              ecs.ComponentsArrayTransaction[transform.PosComponent]
	rotationTransaction         ecs.ComponentsArrayTransaction[transform.RotationComponent]
	sizeTransaction             ecs.ComponentsArrayTransaction[transform.SizeComponent]
	pivotPointTransaction       ecs.ComponentsArrayTransaction[transform.PivotPointComponent]
	parentTransaction           ecs.ComponentsArrayTransaction[transform.ParentComponent]
	parentPivotPointTransaction ecs.ComponentsArrayTransaction[transform.ParentPivotPointComponent]
}

func newTransformTransaction(
	tool tool,
) transform.Transaction {
	return transaction{
		tool,
		tool.posArray.Transaction(),
		tool.rotationArray.Transaction(),
		tool.sizeArray.Transaction(),
		tool.pivotPointArray.Transaction(),
		tool.parentArray.Transaction(),
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
		t.pivotPointTransaction,
		t.parentTransaction,
		t.parentPivotPointTransaction,
	}
}

func (t transaction) Flush() error {
	return ecs.FlushMany(t.Transactions()...)
}
