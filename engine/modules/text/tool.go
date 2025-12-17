package text

import "engine/services/ecs"

type TextTool interface {
	Text() Interface
}

type Interface interface {
	Break() ecs.ComponentsArray[BreakComponent]
	Content() ecs.ComponentsArray[TextComponent]
	Align() ecs.ComponentsArray[TextAlignComponent]
	Color() ecs.ComponentsArray[TextColorComponent]
	FontFamily() ecs.ComponentsArray[FontFamilyComponent]
	FontSize() ecs.ComponentsArray[FontSizeComponent]

	AddDirtySet(ecs.DirtySet)
}
