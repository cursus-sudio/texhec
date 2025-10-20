package texturearray

import (
	"errors"
	"frontend/services/assets"
	"image"
	"shared/services/datastructures"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type TextureArray struct {
	Texture uint32
}

func (r TextureArray) Release() {
	gl.DeleteTextures(1, &r.Texture)
}

func (r TextureArray) Use() {
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, r.Texture)
}

type Factory interface {
	New(images datastructures.SparseArray[uint32, image.Image]) (TextureArray, error)
}

type factory struct {
	assetsStorage assets.AssetsStorage
}

var (
	ErrTexturesHaveToShareSize error = errors.New("all textures have to match size")
)

func (r *factory) New(asset datastructures.SparseArray[uint32, image.Image]) (TextureArray, error) {
	register := TextureArray{}
	images := datastructures.NewSparseArray[uint32, image.Image]()

	bounds := []image.Rectangle{}
	for _, i := range asset.GetIndices() {
		image, _ := asset.Get(i)
		bounds = append(bounds, image.Bounds())
	}

	w, h := 0, 0
	if len(asset.GetValues()) != 0 {
		bounds := asset.GetValues()[0].Bounds()
		w, h = bounds.Dx(), bounds.Dy()
	}

	for _, i := range asset.GetIndices() {
		image, _ := asset.Get(i)

		if w != image.Bounds().Dx() || h != image.Bounds().Dy() {
			return TextureArray{}, ErrTexturesHaveToShareSize
		}

		images.Set(i, image)
	}

	register.Texture = createTexs(w, h, images)

	return register, nil
}
