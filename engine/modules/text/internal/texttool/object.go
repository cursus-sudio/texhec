package texttool

import (
	"engine/modules/text"
	"engine/services/ecs"
)

type object struct {
	transaction
	entity ecs.EntityID

	breakComp  ecs.EntityComponent[text.BreakComponent]
	text       ecs.EntityComponent[text.TextComponent]
	textAlign  ecs.EntityComponent[text.TextAlignComponent]
	textColor  ecs.EntityComponent[text.TextColorComponent]
	fontFamily ecs.EntityComponent[text.FontFamilyComponent]
	fontSize   ecs.EntityComponent[text.FontSizeComponent]
}

func newObject(
	t transaction,
	entity ecs.EntityID,
) text.Object {
	s := object{
		transaction: t,
		entity:      entity,

		breakComp:  t.breakTransaction.GetEntityComponent(entity),
		text:       t.textTransaction.GetEntityComponent(entity),
		textAlign:  t.textAlignTransaction.GetEntityComponent(entity),
		textColor:  t.textColorTransaction.GetEntityComponent(entity),
		fontFamily: t.fontFamilyTransaction.GetEntityComponent(entity),
		fontSize:   t.fontSizeTransaction.GetEntityComponent(entity),
	}
	return s
}

func (o object) Break() ecs.EntityComponent[text.BreakComponent]           { return o.breakComp }
func (o object) Text() ecs.EntityComponent[text.TextComponent]             { return o.text }
func (o object) TextAlign() ecs.EntityComponent[text.TextAlignComponent]   { return o.textAlign }
func (o object) TextColor() ecs.EntityComponent[text.TextColorComponent]   { return o.textColor }
func (o object) FontFamily() ecs.EntityComponent[text.FontFamilyComponent] { return o.fontFamily }
func (o object) FontSize() ecs.EntityComponent[text.FontSizeComponent]     { return o.fontSize }
