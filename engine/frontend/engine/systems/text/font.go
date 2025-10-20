package textsys

import (
	"frontend/engine/components/text"
	"frontend/services/assets"
	"image"
	"image/color"
	"shared/services/datastructures"
	"shared/services/logger"

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
}

func newFontService(
	assets assets.AssetsStorage,
	usedGlyphs datastructures.SparseSet[rune],
	face opentype.FaceOptions,
	logger logger.Logger,
) FontService {
	return &fontService{
		assets,
		usedGlyphs,
		face,
		logger,
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
		width := float32(advance.Ceil()) / 64.
		fontMeta.GlyphsWidth.Set(glyphID, width)

		drawer := font.Drawer{
			Src:  image.NewUniform(color.White),
			Face: fontFace,
		}
		image := getLetterImage(drawer, glyph)
		fontMeta.Images.Set(glyphID, image)
	}

	return fontMeta, nil
}

//	func getTextImage(drawer font.Drawer, text string) *image.RGBA {
//		textBounds, _ := drawer.BoundString(text)
//		// textWidth := textBounds.Max.X - textBounds.Min.X
//		// textHeight := textBounds.Max.Y - textBounds.Min.Y
//
//		drawer.Dot = fixed.Point26_6{
//			X: fixed.I(0) - textBounds.Min.X,
//			Y: fixed.I(0) - textBounds.Min.Y,
//		}
//
//		// rect := image.Rect(0, 0, textWidth.Ceil(), textHeight.Ceil())
//		rect := image.Rect(0, 0, 64, 64)
//		img := image.NewRGBA(rect)
//		drawer.Dst = img
//		drawer.DrawString(text)
//		return img
//	}

const cellSize = 64
const yBaseline = 52

func getLetterImage(drawer font.Drawer, letter rune) *image.RGBA {
	var text string = string(letter)
	textBounds, _ := drawer.BoundString(text)

	rect := image.Rect(0, 0, cellSize, cellSize)
	img := image.NewRGBA(rect)
	drawer.Dst = img

	dotX := fixed.I(0) - textBounds.Min.X
	dotY := fixed.I(yBaseline)

	drawer.Dot = fixed.Point26_6{
		X: dotX,
		Y: dotY,
	}

	drawer.DrawString(text)
	return img
}
