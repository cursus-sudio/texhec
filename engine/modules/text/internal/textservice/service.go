package textservice

import (
	"engine/modules/text"
	"engine/services/ecs"
	"engine/services/logger"
)

type service struct {
	logger logger.Logger

	world ecs.World

	breakArray      ecs.ComponentsArray[text.BreakComponent]
	textArray       ecs.ComponentsArray[text.TextComponent]
	textAlignArray  ecs.ComponentsArray[text.TextAlignComponent]
	textColorArray  ecs.ComponentsArray[text.TextColorComponent]
	fontFamilyArray ecs.ComponentsArray[text.FontFamilyComponent]
	fontSizeArray   ecs.ComponentsArray[text.FontSizeComponent]
}

func NewService(
	w ecs.World,
	logger logger.Logger,
) text.Service {
	return &service{
		logger,
		w,
		ecs.GetComponentsArray[text.BreakComponent](w),
		ecs.GetComponentsArray[text.TextComponent](w),
		ecs.GetComponentsArray[text.TextAlignComponent](w),
		ecs.GetComponentsArray[text.TextColorComponent](w),
		ecs.GetComponentsArray[text.FontFamilyComponent](w),
		ecs.GetComponentsArray[text.FontSizeComponent](w),
	}
}

func (t *service) Break() ecs.ComponentsArray[text.BreakComponent]     { return t.breakArray }
func (t *service) Content() ecs.ComponentsArray[text.TextComponent]    { return t.textArray }
func (t *service) Align() ecs.ComponentsArray[text.TextAlignComponent] { return t.textAlignArray }
func (t *service) Color() ecs.ComponentsArray[text.TextColorComponent] { return t.textColorArray }
func (t *service) FontFamily() ecs.ComponentsArray[text.FontFamilyComponent] {
	return t.fontFamilyArray
}
func (t *service) FontSize() ecs.ComponentsArray[text.FontSizeComponent] { return t.fontSizeArray }

func (t *service) AddDirtySet(set ecs.DirtySet) {
	t.breakArray.AddDirtySet(set)
	t.textAlignArray.AddDirtySet(set)
	t.textColorArray.AddDirtySet(set)
	t.fontFamilyArray.AddDirtySet(set)
	t.fontSizeArray.AddDirtySet(set)
}
