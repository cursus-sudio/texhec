package internal

import (
	"frontend/engine/text"
	"frontend/engine/transform"
	"shared/services/ecs"
	"shared/services/logger"
	"unicode/utf8"

	"github.com/go-gl/mathgl/mgl32"
)

// note: glyph size is (0-1(width is between), 1 (height is const)) and
// in shader its multiplied in shader by font size
type Glyph struct {
	Pos   mgl32.Vec2
	Glyph int32
}

type Layout struct {
	Glyphs   []Glyph
	FontSize uint32
	Font     FontKey
}

type LayoutService interface {
	EntityLayout(ecs.EntityID) (Layout, error)
}

type layoutService struct {
	world           ecs.World
	transformArray  ecs.ComponentsArray[transform.Transform]
	textArray       ecs.ComponentsArray[text.Text]
	fontFamilyArray ecs.ComponentsArray[text.FontFamily]
	fontSizeArray   ecs.ComponentsArray[text.FontSize]
	// overflowArray   ecs.ComponentsArray[text.Overflow]
	breakArray     ecs.ComponentsArray[text.Break]
	textAlignArray ecs.ComponentsArray[text.TextAlign]

	logger      logger.Logger
	fontService FontService
	fontsKeys   FontKeys

	defaultFontFamily text.FontFamily
	defaultFontSize   text.FontSize
	// defaultOverflow   text.Overflow
	defaultBreak     text.Break
	defaultTextAlign text.TextAlign
}

func NewLayoutService(
	world ecs.World,

	logger logger.Logger,
	fontService FontService,
	fontsKeys FontKeys,

	defaultFontFamily text.FontFamily,
	defaultFontSize text.FontSize,
	// defaultOverflow text.Overflow,
	defaultBreak text.Break,
	defaultTextAlign text.TextAlign,
) LayoutService {
	return &layoutService{
		world:           world,
		transformArray:  ecs.GetComponentsArray[transform.Transform](world.Components()),
		textArray:       ecs.GetComponentsArray[text.Text](world.Components()),
		fontFamilyArray: ecs.GetComponentsArray[text.FontFamily](world.Components()),
		fontSizeArray:   ecs.GetComponentsArray[text.FontSize](world.Components()),
		// overflowArray:   ecs.GetComponentsArray[text.Overflow](world.Components()),
		breakArray:     ecs.GetComponentsArray[text.Break](world.Components()),
		textAlignArray: ecs.GetComponentsArray[text.TextAlign](world.Components()),

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

type lineLetter struct {
	letter rune
	xPos   float32
}
type line struct {
	letters []lineLetter
	width   float32
}

// pipeline
// string -> chars
// chars -> lines
// lines -> modify lines (break them where needed)
// line characters -> glyphs

func (s *layoutService) EntityLayout(entity ecs.EntityID) (Layout, error) {
	// TODO add overflow read, text align read and transform modification

	transformComponent, err := s.transformArray.GetComponent(entity)
	if err != nil {
		transformComponent = transform.NewTransform()
		return Layout{}, err
	}
	textComponent, err := s.textArray.GetComponent(entity)
	if err != nil {
		return Layout{}, err
	}
	fontFamily, err := s.fontFamilyArray.GetComponent(entity)
	if err != nil {
		fontFamily = s.defaultFontFamily
	}
	fontSize, err := s.fontSizeArray.GetComponent(entity)
	if err != nil {
		fontSize = s.defaultFontSize
	}
	// overflow, err := s.overflowArray.GetComponent(entity)
	// if err != nil {
	// 	overflow = s.defaultOverflow
	// }
	breakComponent, err := s.breakArray.GetComponent(entity)
	if err != nil {
		breakComponent = s.defaultBreak
	}
	textAlign, err := s.textAlignArray.GetComponent(entity)
	if err != nil {
		textAlign = s.defaultTextAlign
	}

	font, err := s.fontService.AssetFont(fontFamily.FontFamily)
	if err != nil {
		return Layout{}, err
	}

	// create lines letters
	lines := []line{{}}
	lineHeight := 1

	maxWidth := transformComponent.Size.X() / float32(fontSize.FontSize)
	maxHeight := transformComponent.Size.Y() / float32(fontSize.FontSize)

	// generate lines
	var nextLetterIndex int = 0
	for nextLetterIndex < len(textComponent.Text) {
		letter, letterSize := utf8.DecodeRuneInString(textComponent.Text[nextLetterIndex:])
		letterIndex := nextLetterIndex
		nextLetterIndex += letterSize

		if letter == '\n' {
			lines = append(lines, line{})
			continue
		}
		letterLine := lines[len(lines)-1]
		letterWidth, ok := font.GlyphsWidth.Get(uint32(letter))
		if !ok {
			continue
		}

		lineLetter := lineLetter{
			letter: letter,
			xPos:   letterLine.width,
		}

		updatedLine := line{
			letters: append(letterLine.letters, lineLetter),
			width:   letterLine.width + letterWidth,
		}

		var shouldBreak bool = updatedLine.width > maxWidth

		var canBreak bool
		switch breakComponent.Break {
		case text.BreakNone:
			canBreak = false
		case text.BreakAny:
			canBreak = true
		case text.BreakWord:
			canBreak = true
		}

		if !canBreak || !shouldBreak {
			lines[len(lines)-1] = updatedLine
			continue
		}

		var defaultLastLineLetterIndex int = len(updatedLine.letters) - 1
		var lastLineLetterIndex int = defaultLastLineLetterIndex
		switch breakComponent.Break {
		case text.BreakAny:
		case text.BreakNone:
		case text.BreakWord:
			for i := defaultLastLineLetterIndex; i >= 0; i-- {
				if updatedLine.letters[i].letter != ' ' {
					continue
				}
				lastLineLetterIndex = i + 1
				break
			}
		}

		if lastLineLetterIndex < len(updatedLine.letters) {
			updatedLine.width = updatedLine.letters[lastLineLetterIndex].xPos
			updatedLine.letters = updatedLine.letters[:lastLineLetterIndex]
		}
		nextLetterIndex = letterIndex + (lastLineLetterIndex - defaultLastLineLetterIndex)

		lines[len(lines)-1] = updatedLine
		lines = append(lines, line{})
	}

	if len(lines[len(lines)-1].letters) == 0 {
		lines = lines[:len(lines)-1]
	}

	// modify lines
	for _, line := range lines {
		offset := (maxWidth - line.width) * textAlign.Vertical
		for i := 0; i < len(line.letters); i++ {
			line.letters[i].xPos += offset
		}
	}

	var heightOffset float32 = 0
	{
		linesCount := float32(len(lines))
		height := linesCount * float32(lineHeight)
		heightOffset = (maxHeight - height) * textAlign.Horizontal
	}

	// generate glpyhs
	glyphs := []Glyph{}
	for y, line := range lines {
		for _, letter := range line.letters {
			glyph := Glyph{
				Pos: mgl32.Vec2{
					letter.xPos,
					heightOffset + float32(y*int(lineHeight)),
				},
				Glyph: letter.letter,
			}
			glyphs = append(glyphs, glyph)
		}
	}

	layout := Layout{
		Glyphs:   glyphs,
		FontSize: uint32(fontSize.FontSize),
		Font:     s.fontsKeys.GetKey(fontFamily.FontFamily),
	}
	return layout, nil
}
