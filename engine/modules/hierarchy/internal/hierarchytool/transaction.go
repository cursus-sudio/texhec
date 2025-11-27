package hierarchytool

import (
	"engine/modules/hierarchy"
	"engine/services/ecs"
)

type transaction struct {
	tool

	parentTransaction ecs.ComponentsArrayTransaction[hierarchy.ParentComponent]
}

func newTransaction(
	tool tool,
) hierarchy.Transaction {
	return transaction{
		tool,
		tool.parentArray.Transaction(),
	}
}

func (t transaction) GetObject(entity ecs.EntityID) hierarchy.Object {
	return newObject(t, entity)
}

func (t transaction) Transactions() []ecs.AnyComponentsArrayTransaction {
	return []ecs.AnyComponentsArrayTransaction{t.parentTransaction}
}

func (t transaction) Flush() error {
	return ecs.FlushMany(t.Transactions()...)
}
