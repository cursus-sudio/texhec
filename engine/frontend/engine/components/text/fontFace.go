package text

import (
	"golang.org/x/image/font/opentype"
)

type FontFaceAsset interface {
	Font() *opentype.Font
	Release()
}

type fontFaceAsset struct {
	font opentype.Font
}

func (face fontFaceAsset) Font() *opentype.Font { return &face.font }

func (face fontFaceAsset) Release() {}

func NewFontFaceAsset(raw opentype.Font) FontFaceAsset {
	return fontFaceAsset{raw}
}
