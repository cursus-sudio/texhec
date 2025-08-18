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

type TextureStorageAsset interface {
	Image() image.Image
}

type textureStorageAsset struct {
	image image.Image
}

func NewTextureStorageAsset(
	image image.Image,
) TextureStorageAsset {
	return &textureStorageAsset{
		image: image,
	}
}

func (a *textureStorageAsset) Image() image.Image { return a.image }
