package texturearray

import (
	"errors"
	"frontend/engine/components/texture"
	"frontend/services/assets"
	"image"
	"shared/services/datastructures"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type TextureArray struct {
	Assets  datastructures.Set[assets.AssetID]
	Texture uint32
}

func (r TextureArray) Release() {
	gl.DeleteTextures(1, &r.Texture)
}

func (r TextureArray) Use() {
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, r.Texture)
}

type Factory interface {
	Add(asset ...assets.AssetID)
	New() (TextureArray, error)
}

type factory struct {
	assetsStorage assets.AssetsStorage
	assets        []assets.AssetID
}

var (
	ErrTexturesHaveToShareSize error = errors.New("all textures have to match size")
)

func (r *factory) Add(asset ...assets.AssetID) {
	r.assets = append(r.assets, asset...)
}

func (r *factory) New() (TextureArray, error) {
	register := TextureArray{
		Assets: datastructures.NewSet[assets.AssetID](),
	}
	images := []image.Image{}

	w, h := 0, 0
	for i, assetID := range r.assets {
		asset, err := assets.StorageGet[texture.TextureAsset](r.assetsStorage, assetID)
		if err != nil {
			err := errors.Join(
				err,
				errors.New("creating world texture register"),
			)
			return TextureArray{}, err
		}

		image := asset.Image()
		if i == 0 {
			w, h = image.Bounds().Dx(), image.Bounds().Dy()
		}
		if w != image.Bounds().Dx() || h != image.Bounds().Dy() {
			return TextureArray{}, ErrTexturesHaveToShareSize
		}

		images = append(images, image)
		register.Assets.Add(assetID)
	}

	register.Texture = createTexs(w, h, images)

	return register, nil
}
