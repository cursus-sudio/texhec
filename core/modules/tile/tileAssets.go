package tile

import (
	"engine/modules/render"
	gtexture "engine/services/graphics/texture"
	"image"
)

//

type BiomAsset interface {
	Images() [15][]image.Image
	Res() image.Rectangle
	AspectRatio() image.Rectangle
}

type biomAsset struct {
	images      [15][]image.Image
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

// src images should be:
// - 1111
// - 1110
// - 1010
// - 1001
// - 0001
func NewBiomAsset(srcImages [6][]image.Image) (BiomAsset, error) {
	var res image.Rectangle
	for i, imgages := range srcImages {
		for _, img := range imgages {
			bounds := img.Bounds()
			bounds = image.Rect(0, 0, bounds.Dx(), bounds.Dy())
			if i == 0 {
				res = bounds
				continue
			}
			if res != bounds {
				return nil, render.ErrTextureAssetImagesHasToMatchResolution
			}
		}
	}

	aspectRatio := image.Rect(0, 0, res.Dx(), res.Dy())
	divisor := greatestCommonDivisor(aspectRatio.Max.X, aspectRatio.Max.Y)

	aspectRatio.Max.X /= divisor
	aspectRatio.Max.Y /= divisor

	// example rotate 0123 -> 2031 -> 3210 -> 1302
	// example rotate 1000 -> 0100 -> 0001 -> 0010
	// src images -> rotated:
	// - 1111
	// - 1110 -> 1101 -> 0111 -> 1011
	// - 1010 -> 1100 -> 0101 -> 0011
	// - 1001 -> 0110
	// - 0001 -> 0010 -> 1000 -> 0100

	// this is new version of tiles
	// it uses one more tile to define biom
	// but it maintains horizontal and vertical axis
	dstImages := [15][]image.Image{}
	// this code either can be organized by srcImages or dstImages
	// ordering in ascending order for srcImages sounds best for maintanance
	// but due to additional for loops dstImages in ascending order are choosen
	for _, img := range srcImages[0] {
		dstImages[2] = append(dstImages[2], img)
		dstImages[11] = append(dstImages[11], gtexture.NewImage(img).FlipV().Image())
	}
	dstImages[14] = append(dstImages[14], srcImages[1]...)
	for _, img := range srcImages[2] {
		dstImages[6] = append(dstImages[6], gtexture.NewImage(img).FlipV().Image())
		dstImages[10] = append(dstImages[10], gtexture.NewImage(img).FlipHV().Image())
		dstImages[12] = append(dstImages[12], img)
		dstImages[13] = append(dstImages[13], gtexture.NewImage(img).FlipH().Image())
	}
	for _, img := range srcImages[3] {
		dstImages[4] = append(dstImages[4], img)
		dstImages[9] = append(dstImages[9], gtexture.NewImage(img).FlipH().Image())
	}
	for _, img := range srcImages[4] {
		dstImages[5] = append(dstImages[5], img)
		dstImages[8] = append(dstImages[8], gtexture.NewImage(img).FlipH().Image())
	}
	for _, img := range srcImages[5] {
		dstImages[0] = append(dstImages[0], gtexture.NewImage(img).FlipH().Image())
		dstImages[1] = append(dstImages[1], img)
		dstImages[3] = append(dstImages[3], gtexture.NewImage(img).FlipHV().Image())
		dstImages[7] = append(dstImages[7], gtexture.NewImage(img).FlipV().Image())
	}

	asset := &biomAsset{
		images:      dstImages,
		res:         res,
		aspectRatio: aspectRatio,
	}
	return asset, nil
}

func (a *biomAsset) Images() [15][]image.Image    { return a.images }
func (a *biomAsset) Res() image.Rectangle         { return a.res }
func (a *biomAsset) AspectRatio() image.Rectangle { return a.aspectRatio }
func (a *biomAsset) Release()                     {}
