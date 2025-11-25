package textrenderer

import (
	"engine/modules/text"
	"engine/services/assets"
	"engine/services/datastructures"
	"engine/services/graphics/texture"
	"engine/services/logger"
	"image"
	"image/color"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type FontKey uint32

type Font struct {
	GlyphsWidth datastructures.SparseArray[uint32, float32]
	Images      datastructures.SparseArray[uint32, image.Image]
}

type FontService interface {
	AssetFont(assets.AssetID) (Font, error)
}

type fontService struct {
	assets      assets.AssetsStorage
	usedGlyphs  datastructures.SparseSet[rune]
	faceOptions opentype.FaceOptions
	logger      logger.Logger

	cellSize, yBaseline int
}

func NewFontService(
	assets assets.AssetsStorage,
	usedGlyphs datastructures.SparseSet[rune],
	face opentype.FaceOptions,
	logger logger.Logger,
	cellSize, yBaseline int,
) FontService {
	return &fontService{
		assets,
		usedGlyphs,
		face,
		logger,
		cellSize,
		yBaseline,
	}
}

func (s *fontService) AssetFont(assetID assets.AssetID) (Font, error) {
	fontMeta := Font{
		GlyphsWidth: datastructures.NewSparseArray[uint32, float32](),
		Images:      datastructures.NewSparseArray[uint32, image.Image](),
	}
	asset, err := assets.StorageGet[text.FontFaceAsset](s.assets, assetID)
	if err != nil {
		return Font{}, err
	}
	face := asset.Font()

	fontFace, err := opentype.NewFace(face, &s.faceOptions)
	if err != nil {
		return Font{}, err
	}

	glyphs := s.usedGlyphs.GetIndices()
	for _, glyph := range glyphs {
		glyphID := uint32(glyph)
		_, advance, _ := fontFace.GlyphBounds(glyph)
		width := float32(advance.Ceil()) / float32(s.faceOptions.Size)
		fontMeta.GlyphsWidth.Set(glyphID, width)

		drawer := font.Drawer{
			Src:  image.NewUniform(color.White),
			Face: fontFace,
		}
		image := s.getLetterImage(drawer, glyph)
		fontMeta.Images.Set(glyphID, texture.FlipImage(image))
	}

	return fontMeta, nil
}

func (s *fontService) getLetterImage(drawer font.Drawer, letter rune) *image.RGBA {
	var text string = string(letter)
	textBounds, _ := drawer.BoundString(text)

	rect := image.Rect(0, 0, s.cellSize, s.cellSize)
	img := image.NewRGBA(rect)
	drawer.Dst = img

	dotX := fixed.I(0) - textBounds.Min.X
	dotY := fixed.I(s.yBaseline)

	drawer.Dot = fixed.Point26_6{
		X: dotX,
		Y: dotY,
	}

	drawer.DrawString(text)
	return img
}
