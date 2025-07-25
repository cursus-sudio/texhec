package texture

import (
	"frontend/services/assets"
	"frontend/services/graphics/texture"
	"io"
)

type Texture struct {
	ID assets.AssetID
}

func NewTexture(id assets.AssetID) Texture {
	return Texture{ID: id}
}

type TextureStorageAsset interface {
	assets.StorageAsset
	Reader() io.Reader
}

type textureStorageAsset struct {
	reader io.Reader
}

func NewTextureStorageAsset(
	reader io.Reader,
) TextureStorageAsset {
	return &textureStorageAsset{
		reader: reader,
	}
}

func (a *textureStorageAsset) Reader() io.Reader { return a.reader }

func (a *textureStorageAsset) Cache() (assets.CachedAsset, error) {
	t, err := texture.NewTexture(a.reader)
	if err != nil {
		return nil, err
	}
	return &textureCachedAsset{texture: t}, nil
}

//

type TextureCachedAsset interface {
	assets.CachedAsset
	Texture() texture.Texture
}

type textureCachedAsset struct {
	texture texture.Texture
}

func (asset *textureCachedAsset) Texture() texture.Texture { return asset.texture }
func (asset *textureCachedAsset) Release()                 { asset.texture.Release() }
