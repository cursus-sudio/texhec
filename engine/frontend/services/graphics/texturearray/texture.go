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
	Add(asset datastructures.SparseArray[uint32, assets.AssetID])
	New() (TextureArray, error)
}

type factory struct {
	assetsStorage assets.AssetsStorage
	assets        datastructures.SparseArray[uint32, assets.AssetID]
}

var (
	ErrTexturesHaveToShareSize error = errors.New("all textures have to match size")
)

func (r *factory) Add(asset datastructures.SparseArray[uint32, assets.AssetID]) {
	for _, index := range asset.GetIndices() {
		value, _ := asset.Get(index)
		r.assets.Set(index, value)
	}
}

func (r *factory) New() (TextureArray, error) {
	register := TextureArray{
		Assets: datastructures.NewSet[assets.AssetID](),
	}
	images := datastructures.NewSparseArray[uint32, image.Image]()

	w, h := 0, 0
	for _, i := range r.assets.GetIndices() {
		assetID, _ := r.assets.Get(i)
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

		images.Set(i, image)
		register.Assets.Add(assetID)
	}

	register.Texture = createTexs(w, h, images)

	return register, nil
}
