package texttool

import (
	"engine/modules/text"
	"engine/services/ecs"
	"engine/services/logger"
)

type tool struct {
	logger logger.Logger

	world ecs.World

	breakArray      ecs.ComponentsArray[text.BreakComponent]
	textArray       ecs.ComponentsArray[text.TextComponent]
	textAlignArray  ecs.ComponentsArray[text.TextAlignComponent]
	textColorArray  ecs.ComponentsArray[text.TextColorComponent]
	fontFamilyArray ecs.ComponentsArray[text.FontFamilyComponent]
	fontSizeArray   ecs.ComponentsArray[text.FontSizeComponent]
}

func NewTool(
	logger logger.Logger,
) ecs.ToolFactory[text.TextTool] {
	return ecs.NewToolFactory(func(w ecs.World) text.TextTool {
		return tool{
			logger,
			w,
			ecs.GetComponentsArray[text.BreakComponent](w),
			ecs.GetComponentsArray[text.TextComponent](w),
			ecs.GetComponentsArray[text.TextAlignComponent](w),
			ecs.GetComponentsArray[text.TextColorComponent](w),
			ecs.GetComponentsArray[text.FontFamilyComponent](w),
			ecs.GetComponentsArray[text.FontSizeComponent](w),
		}
	})
}

func (t tool) Text() text.Interface { return t }

func (t tool) Break() ecs.ComponentsArray[text.BreakComponent]           { return t.breakArray }
func (t tool) Content() ecs.ComponentsArray[text.TextComponent]          { return t.textArray }
func (t tool) TextAlign() ecs.ComponentsArray[text.TextAlignComponent]   { return t.textAlignArray }
func (t tool) TextColor() ecs.ComponentsArray[text.TextColorComponent]   { return t.textColorArray }
func (t tool) FontFamily() ecs.ComponentsArray[text.FontFamilyComponent] { return t.fontFamilyArray }
func (t tool) FontSize() ecs.ComponentsArray[text.FontSizeComponent]     { return t.fontSizeArray }

func (t tool) AddDirtySet(set ecs.DirtySet) {
	t.breakArray.AddDirtySet(set)
	t.textAlignArray.AddDirtySet(set)
	t.textColorArray.AddDirtySet(set)
	t.fontFamilyArray.AddDirtySet(set)
	t.fontSizeArray.AddDirtySet(set)
}
