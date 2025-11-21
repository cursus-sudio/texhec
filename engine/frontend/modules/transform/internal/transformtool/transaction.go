package transformtool

import (
	"frontend/modules/transform"
	"shared/services/ecs"
)

type transformTransaction struct {
	transformTool

	posTransaction              ecs.ComponentsArrayTransaction[transform.PosComponent]
	rotationTransaction         ecs.ComponentsArrayTransaction[transform.RotationComponent]
	sizeTransaction             ecs.ComponentsArrayTransaction[transform.SizeComponent]
	pivotPointTransaction       ecs.ComponentsArrayTransaction[transform.PivotPointComponent]
	parentTransaction           ecs.ComponentsArrayTransaction[transform.ParentComponent]
	parentPivotPointTransaction ecs.ComponentsArrayTransaction[transform.ParentPivotPointComponent]
}

func newTransformTransaction(
	tool transformTool,
) transform.TransformTransaction {
	return transformTransaction{
		tool,
		tool.posArray.Transaction(),
		tool.rotationArray.Transaction(),
		tool.sizeArray.Transaction(),
		tool.pivotPointArray.Transaction(),
		tool.parentArray.Transaction(),
		tool.parentPivotPointArray.Transaction(),
	}
}

func (t transformTransaction) GetEntity(entity ecs.EntityID) transform.EntityTransform {
	return newEntityTransform(t, entity)
}

func (t transformTransaction) Transactions() []ecs.AnyComponentsArrayTransaction {
	return []ecs.AnyComponentsArrayTransaction{
		t.posTransaction,
		t.rotationTransaction,
		t.sizeTransaction,
		t.pivotPointTransaction,
		t.parentTransaction,
		t.parentPivotPointTransaction,
	}
}

func (t transformTransaction) Flush() error {
	return ecs.FlushMany(t.Transactions()...)
}
