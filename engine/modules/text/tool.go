package text

import "engine/services/ecs"

type Tool interface {
	Transaction() Transaction
	Query(ecs.LiveQueryBuilder) ecs.LiveQueryBuilder
}

type Transaction interface {
	GetObject(ecs.EntityID) Object
	Transactions() []ecs.AnyComponentsArrayTransaction
	Flush() error
}

type Object interface {
	Break() ecs.EntityComponent[BreakComponent]
	Text() ecs.EntityComponent[TextComponent]
	TextAlign() ecs.EntityComponent[TextAlignComponent]
	TextColor() ecs.EntityComponent[TextColorComponent]
	FontFamily() ecs.EntityComponent[FontFamilyComponent]
	FontSize() ecs.EntityComponent[FontSizeComponent]
}
