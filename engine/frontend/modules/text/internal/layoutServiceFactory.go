package internal

import (
	"frontend/modules/text"
	"shared/services/ecs"
	"shared/services/logger"
)

type LayoutServiceFactory interface {
	New(ecs.World) LayoutService
}

type layoutServiceFactory struct {
	logger      logger.Logger
	fontService FontService
	fontsKeys   FontKeys

	defaultFontFamily text.FontFamilyComponent
	defaultFontSize   text.FontSizeComponent
	// defaultOverflow   text.Overflow
	defaultBreak     text.BreakComponent
	defaultTextAlign text.TextAlignComponent
}

func NewLayoutServiceFactory(
	logger logger.Logger,
	fontService FontService,
	fontsKeys FontKeys,

	defaultFontFamily text.FontFamilyComponent,
	defaultFontSize text.FontSizeComponent,
	// defaultOverflow text.Overflow,
	defaultBreak text.BreakComponent,
	defaultTextAlign text.TextAlignComponent,
) LayoutServiceFactory {
	return &layoutServiceFactory{
		logger:      logger,
		fontService: fontService,
		fontsKeys:   fontsKeys,

		defaultFontFamily: defaultFontFamily,
		defaultFontSize:   defaultFontSize,
		// defaultOverflow:   defaultOverflow,
		defaultBreak:     defaultBreak,
		defaultTextAlign: defaultTextAlign,
	}
}

func (f *layoutServiceFactory) New(world ecs.World) LayoutService {
	return NewLayoutService(
		world,
		f.logger,
		f.fontService,
		f.fontsKeys,
		f.defaultFontFamily,
		f.defaultFontSize,
		// f.defaultOverflow,
		f.defaultBreak,
		f.defaultTextAlign,
	)
}
