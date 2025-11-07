package render

import (
	"frontend/services/assets"
	"image"
)

type TextureComponent struct {
	Asset assets.AssetID
	Frame int
}

func NewTexture(asset assets.AssetID) TextureComponent {
	return TextureComponent{Asset: asset}
}

func (c TextureComponent) SetFrame(frame int) TextureComponent {
	c.Frame = frame
	return c
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
