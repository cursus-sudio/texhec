package text

import "engine/services/ecs"

type Text interface {
	Text() Interface
}

type Interface interface {
	Break() ecs.ComponentsArray[BreakComponent]
	TextContent() ecs.ComponentsArray[TextComponent]
	TextAlign() ecs.ComponentsArray[TextAlignComponent]
	TextColor() ecs.ComponentsArray[TextColorComponent]
	FontFamily() ecs.ComponentsArray[FontFamilyComponent]
	FontSize() ecs.ComponentsArray[FontSizeComponent]

	AddDirtySet(ecs.DirtySet)
}
