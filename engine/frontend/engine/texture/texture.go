package texture

import (
	"frontend/services/assets"
	"image"
)

type Texture struct {
	ID assets.AssetID
}

func NewTexture(id assets.AssetID) Texture {
	return Texture{ID: id}
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
