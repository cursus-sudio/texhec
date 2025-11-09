package render

import (
	"frontend/services/assets"
	"image"
)

type TextureComponent struct {
	Asset assets.AssetID
}

func NewTexture(asset assets.AssetID) TextureComponent {
	return TextureComponent{Asset: asset}
}

//

// frame is normalized
type TextureFrameComponent struct {
	FrameNormalized float64
}

func NewTextureFrameComponent(frameNormalized float64) TextureFrameComponent {
	return TextureFrameComponent{
		FrameNormalized: max(min(frameNormalized, 1), 0),
	}
}

func DefaultTextureFrameComponent() TextureFrameComponent {
	return TextureFrameComponent{0}
}

func (c TextureFrameComponent) GetFrame(frameLen int) int {
	return min(
		int(c.FrameNormalized*float64(frameLen)),
		frameLen-1,
	)
}

func (c1 TextureFrameComponent) Blend(c2 TextureFrameComponent, mix64 float64) TextureFrameComponent {
	invMix64 := 1.0 - mix64
	frame := c1.FrameNormalized*invMix64 + c2.FrameNormalized*mix64
	return TextureFrameComponent{FrameNormalized: frame}
}

//

type TextureAsset interface {
	Images() []image.Image
}

type textureAsset struct {
	images []image.Image
}

func NewTextureStorageAsset(
	images ...image.Image,
) TextureAsset {
	return &textureAsset{
		images: images,
	}
}

func (a *textureAsset) Images() []image.Image { return a.images }
