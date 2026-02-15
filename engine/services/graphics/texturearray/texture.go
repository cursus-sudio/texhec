package texturearray

import (
	"engine/services/datastructures"
	"errors"
	"image"

	"github.com/go-gl/gl/v4.5-core/gl"
)

type TextureArray struct {
	Texture     uint32
	ImagesCount int
}

func (r TextureArray) Release() {
	gl.DeleteTextures(1, &r.Texture)
}

func (r TextureArray) Bind() {
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, r.Texture)
}

type Factory interface {
	New(images datastructures.SparseArray[uint32, image.Image]) (TextureArray, error)
	NewFromSlice([]image.Image) (TextureArray, error)
	Wrap(wrapper func(TextureArray))
}

type factory struct {
	wrappers []func(TextureArray)
}

var (
	ErrTexturesHaveToShareSize error = errors.New("all textures have to match size")
)

func (f *factory) New(asset datastructures.SparseArray[uint32, image.Image]) (TextureArray, error) {
	textureArray := TextureArray{}
	images := datastructures.NewSparseArray[uint32, image.Image]()

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

	textureArray.Texture = createTexs(w, h, images)
	textureArray.ImagesCount = images.Size()

	for _, wrapper := range f.wrappers {
		wrapper(textureArray)
	}

	return textureArray, nil
}

func (f *factory) NewFromSlice(images []image.Image) (TextureArray, error) {
	arr := datastructures.NewSparseArray[uint32, image.Image]()
	for i, image := range images {
		arr.Set(uint32(i), image)
	}
	return f.New(arr)
}

func (f *factory) Wrap(wrapper func(TextureArray)) {
	f.wrappers = append(f.wrappers, wrapper)
}
