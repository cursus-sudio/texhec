package textrenderer

import (
	"engine/modules/assets"
	"engine/modules/text"
	"engine/services/datastructures"
	"engine/services/logger"

	"golang.org/x/image/font/opentype"
)

type FontKey uint32

type FontService interface {
	AssetFont(assets.ID) (text.Glyphs, error)
}

type fontService struct {
	assets      assets.Service
	usedGlyphs  datastructures.SparseSet[rune]
	faceOptions opentype.FaceOptions
	logger      logger.Logger

	cellSize, yBaseline int
}

func NewFontService(
	assets assets.Service,
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

// temporary fix for performance
// there should be in public font asset type and added in factory

//

func (s *fontService) AssetFont(assetID assets.ID) (text.Glyphs, error) {
	asset, err := assets.GetAsset[text.FontFaceAsset](s.assets, assetID)
	if err != nil {
		return text.Glyphs{}, err
	}
	return asset.Glyphs(), nil
}
