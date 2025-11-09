package render

import (
	"errors"
	"frontend/services/assets"
	"image"
	"math"
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

var (
	ErrCannotBlendTextureBetweenDifferentAssets = errors.New("cannot blend texture components between different assets")
)

func (c1 TextureComponent) Blend(c2 TextureComponent, mix32 float32) (TextureComponent, error) {
	if c1.Asset != c2.Asset {
		return TextureComponent{}, ErrCannotBlendTextureBetweenDifferentAssets
	}
	invMix32 := 1.0 - mix32
	frame := float32(c1.Frame)*invMix32 + float32(c2.Frame)*mix32

	return TextureComponent{
		Asset: c1.Asset,
		Frame: int(math.Round(float64(frame))),
	}, nil
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
