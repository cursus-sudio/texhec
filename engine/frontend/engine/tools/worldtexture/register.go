package worldtexture

import (
	"errors"
	"frontend/engine/components/texture"
	"frontend/services/assets"
	"frontend/services/datastructures"
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type WorldTextureRegister struct {
	Assets  datastructures.Set[assets.AssetID]
	Texture uint32
}

func (r WorldTextureRegister) Release() {
	gl.DeleteTextures(1, &r.Texture)
}

func (r WorldTextureRegister) Use() {
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, r.Texture)
}

type RegisterFactory interface {
	New(assets ...assets.AssetID) (WorldTextureRegister, error)
}

type registerFactory struct {
	assetsStorage assets.AssetsStorage
}

var (
	ErrTexturesHaveToShareSize error = errors.New("all textures have to match size")
)

func (r *registerFactory) New(textureAssets ...assets.AssetID) (WorldTextureRegister, error) {
	register := WorldTextureRegister{
		Assets: datastructures.NewSet[assets.AssetID](),
	}
	images := []image.Image{}

	w, h := 0, 0
	for i, assetID := range textureAssets {
		asset, err := assets.StorageGet[texture.TextureAsset](r.assetsStorage, assetID)
		if err != nil {
			err := errors.Join(
				err,
				errors.New("creating world texture register"),
			)
			return WorldTextureRegister{}, err
		}

		image := asset.Image()
		if i == 0 {
			w, h = image.Bounds().Dx(), image.Bounds().Dy()
		}
		if w != image.Bounds().Dx() || h != image.Bounds().Dy() {
			return WorldTextureRegister{}, ErrTexturesHaveToShareSize
		}

		images = append(images, image)
		register.Assets.Add(assetID)
	}

	register.Texture = CreateTexs(w, h, images)

	return register, nil
}
