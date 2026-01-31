package render

import (
	"engine/services/assets"
	"errors"
	"image"
)

type TextureComponent struct {
	Asset assets.AssetID
}

func NewTexture(asset assets.AssetID) TextureComponent {
	return TextureComponent{Asset: asset}
}

//

// frame is normalized
type TextureFrameComponent struct {
	FrameNormalized float64
}

func NewTextureFrame(frameNormalized float64) TextureFrameComponent {
	return TextureFrameComponent{
		FrameNormalized: max(min(frameNormalized, 1), 0),
	}
}

func (c TextureFrameComponent) GetFrame(frameLen int) int {
	return min(
		int(c.FrameNormalized*float64(frameLen)),
		frameLen-1,
	)
}

func (c1 TextureFrameComponent) Lerp(c2 TextureFrameComponent, mix32 float32) TextureFrameComponent {
	mix64 := float64(mix32)
	return TextureFrameComponent{c1.FrameNormalized*(1-mix64) + c2.FrameNormalized*mix64}
}

//

var (
	ErrTextureAssetRequiresImages             error = errors.New("texture asset requires images")
	ErrTextureAssetImagesHasToMatchResolution error = errors.New("images have to have the same resolution")
)

type TextureAsset interface {
	Images() []image.Image
	Res() image.Rectangle
	AspectRatio() image.Rectangle
}

type textureAsset struct {
	images      []image.Image
	res         image.Rectangle
	aspectRatio image.Rectangle
}

// greatestCommonDivisor calculates the Greatest Common Divisor of two integers (a and b)
func greatestCommonDivisor(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func NewTextureAsset(
	images ...image.Image,
) (TextureAsset, error) {
	if len(images) == 0 {
		return nil, ErrTextureAssetRequiresImages
	}
	var res image.Rectangle
	for i, img := range images {
		bounds := img.Bounds()
		bounds = image.Rect(0, 0, bounds.Dx(), bounds.Dy())
		if i == 0 {
			res = bounds
			continue
		}
		if res != bounds {
			return nil, ErrTextureAssetImagesHasToMatchResolution
		}
	}

	aspectRatio := image.Rect(0, 0, res.Dx(), res.Dy())
	divisor := greatestCommonDivisor(aspectRatio.Max.X, aspectRatio.Max.Y)

	aspectRatio.Max.X /= divisor
	aspectRatio.Max.Y /= divisor

	asset := &textureAsset{
		images:      images,
		res:         res,
		aspectRatio: aspectRatio,
	}
	return asset, nil
}

func (a *textureAsset) Images() []image.Image        { return a.images }
func (a *textureAsset) Res() image.Rectangle         { return a.res }
func (a *textureAsset) AspectRatio() image.Rectangle { return a.aspectRatio }
func (a *textureAsset) Release()                     {}
