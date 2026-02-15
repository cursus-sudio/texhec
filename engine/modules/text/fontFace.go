package text

import (
	"engine/services/datastructures"
	"image"

	"golang.org/x/image/font/opentype"
)

type Glyphs struct {
	GlyphsWidth datastructures.SparseArray[uint32, float32]
	Images      datastructures.SparseArray[uint32, image.Image]
}

type FontFaceAsset interface {
	Font() *opentype.Font
	Glyphs() Glyphs
	Release()
}

type fontAsset struct {
	font   opentype.Font
	glyphs Glyphs
}

func (f fontAsset) Font() *opentype.Font { return &f.font }
func (f fontAsset) Glyphs() Glyphs       { return f.glyphs }
func (f fontAsset) Release()             {}

func NewFontAsset(raw opentype.Font, glyphs Glyphs) FontFaceAsset {
	return fontAsset{raw, glyphs}
}
