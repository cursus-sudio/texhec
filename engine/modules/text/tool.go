package text

import (
	"engine/modules/camera"
	"engine/modules/groups"
	"engine/modules/transform"
	"engine/services/ecs"
)

type ToolFactory ecs.ToolFactory[World, TextTool]
type TextTool interface {
	Text() Interface
}
type World interface {
	ecs.World
	groups.GroupsTool
	camera.CameraTool
	transform.TransformTool
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
