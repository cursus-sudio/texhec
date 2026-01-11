package textrenderer

import (
	"engine/modules/text"
	"engine/modules/transform"
	"engine/services/ecs"
	"engine/services/logger"
	"unicode/utf8"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ogiusek/ioc/v2"
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
	Logger logger.Logger `inject:"1"`

	World     ecs.World         `inject:"1"`
	Transform transform.Service `inject:"1"`
	Text      text.Service      `inject:"1"`

	FontService FontService `inject:"1"`
	FontsKeys   FontKeys    `inject:"1"`

	defaultFontFamily text.FontFamilyComponent
	defaultFontSize   text.FontSizeComponent
	// defaultOverflow   text.Overflow
	defaultBreak     text.BreakComponent
	defaultTextAlign text.TextAlignComponent
}

func NewLayoutService(c ioc.Dic,
	defaultFontFamily text.FontFamilyComponent,
	defaultFontSize text.FontSizeComponent,
	// defaultOverflow text.Overflow,
	defaultBreak text.BreakComponent,
	defaultTextAlign text.TextAlignComponent,
) LayoutService {
	s := ioc.GetServices[*layoutService](c)

	s.defaultFontFamily = defaultFontFamily
	s.defaultFontSize = defaultFontSize
	// s.defaultOverflow = defaultOverflow
	s.defaultBreak = defaultBreak
	s.defaultTextAlign = defaultTextAlign

	return s
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

	size, _ := s.Transform.AbsoluteSize().Get(entity)
	textComponent, ok := s.Text.Content().Get(entity)
	if !ok {
		return Layout{}, nil
	}
	fontFamily, ok := s.Text.FontFamily().Get(entity)
	if !ok {
		fontFamily = s.defaultFontFamily
	}
	fontSize, ok := s.Text.FontSize().Get(entity)
	if !ok {
		fontSize = s.defaultFontSize
	}
	// overflow, err := s.overflowArray.GetComponent(entity)
	// if err != nil {
	// 	overflow = s.defaultOverflow
	// }
	breakComponent, ok := s.Text.Break().Get(entity)
	if !ok {
		breakComponent = s.defaultBreak
	}
	textAlign, ok := s.Text.Align().Get(entity)
	if !ok {
		textAlign = s.defaultTextAlign
	}

	font, err := s.FontService.AssetFont(fontFamily.FontFamily)
	if err != nil {
		return Layout{}, err
	}

	// create lines letters
	lines := []line{{}}
	lineHeight := 1

	maxWidth := size.Size.X() / float32(fontSize.FontSize)
	maxHeight := size.Size.Y() / float32(fontSize.FontSize)

	// generate lines
	var nextLetterIndex = 0
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

		var shouldBreak = updatedLine.width > maxWidth

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

		var defaultLastLineLetterIndex = len(updatedLine.letters) - 1
		var lastLineLetterIndex = defaultLastLineLetterIndex
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
		lastLineLetterIndex = max(1, lastLineLetterIndex)

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

	var heightOffset float32
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
		Font:     s.FontsKeys.GetKey(fontFamily.FontFamily),
	}
	return layout, nil
}
