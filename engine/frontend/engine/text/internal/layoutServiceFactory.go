package internal

import (
	"frontend/engine/text"
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

	defaultFontFamily text.FontFamily
	defaultFontSize   text.FontSize
	// defaultOverflow   text.Overflow
	defaultBreak     text.Break
	defaultTextAlign text.TextAlign
}

func NewLayoutServiceFactory(
	logger logger.Logger,
	fontService FontService,
	fontsKeys FontKeys,

	defaultFontFamily text.FontFamily,
	defaultFontSize text.FontSize,
	// defaultOverflow text.Overflow,
	defaultBreak text.Break,
	defaultTextAlign text.TextAlign,
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
