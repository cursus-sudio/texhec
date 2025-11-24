package texttool

import (
	"frontend/modules/text"
	"shared/services/ecs"
	"shared/services/logger"
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
) ecs.ToolFactory[text.Tool] {
	return ecs.NewToolFactory(func(w ecs.World) text.Tool {
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

func (tool tool) Transaction() text.Transaction {
	return newTransaction(tool)
}

func (tool tool) Query(b ecs.LiveQueryBuilder) ecs.LiveQueryBuilder {
	return b.Require(ecs.GetComponentType(text.TextComponent{})).Track(
		ecs.GetComponentType(text.BreakComponent{}),
		ecs.GetComponentType(text.TextAlignComponent{}),
		ecs.GetComponentType(text.TextColorComponent{}),
		ecs.GetComponentType(text.FontFamilyComponent{}),
		ecs.GetComponentType(text.FontSizeComponent{}),
	)
}
