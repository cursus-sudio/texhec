package texture

import (
	"frontend/services/assets"
	"image"
)

type TextureComponent struct {
	ID assets.AssetID
}

func NewTexture(id assets.AssetID) TextureComponent {
	return TextureComponent{ID: id}
}

type TextureAsset interface {
	Image() image.Image
}

type textureAsset struct {
	image image.Image
}

func NewTextureStorageAsset(
	image image.Image,
) TextureAsset {
	return &textureAsset{
		image: image,
	}
}

func (a *textureAsset) Image() image.Image { return a.image }
