package texttool

import (
	"frontend/modules/text"
	"shared/services/ecs"
)

type transaction struct {
	tool

	breakTransaction      ecs.ComponentsArrayTransaction[text.BreakComponent]
	textTransaction       ecs.ComponentsArrayTransaction[text.TextComponent]
	textAlignTransaction  ecs.ComponentsArrayTransaction[text.TextAlignComponent]
	textColorTransaction  ecs.ComponentsArrayTransaction[text.TextColorComponent]
	fontFamilyTransaction ecs.ComponentsArrayTransaction[text.FontFamilyComponent]
	fontSizeTransaction   ecs.ComponentsArrayTransaction[text.FontSizeComponent]
}

func newTransaction(
	tool tool,
) text.Transaction {
	return transaction{
		tool,
		tool.breakArray.Transaction(),
		tool.textArray.Transaction(),
		tool.textAlignArray.Transaction(),
		tool.textColorArray.Transaction(),
		tool.fontFamilyArray.Transaction(),
		tool.fontSizeArray.Transaction(),
	}
}

func (t transaction) GetObject(entity ecs.EntityID) text.Object {
	return newObject(t, entity)
}

func (t transaction) Transactions() []ecs.AnyComponentsArrayTransaction {
	return []ecs.AnyComponentsArrayTransaction{
		t.breakTransaction,
		t.textTransaction,
		t.textAlignTransaction,
		t.textColorTransaction,
		t.fontFamilyTransaction,
		t.fontSizeTransaction,
	}
}

func (t transaction) Flush() error {
	return ecs.FlushMany(t.Transactions()...)
}
