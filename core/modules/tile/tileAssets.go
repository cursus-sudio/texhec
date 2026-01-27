package tile

import (
	"engine/modules/render"
	gtexture "engine/services/graphics/texture"
	"image"
)

//

type BiomAsset interface {
	Images() [15]image.Image
	Res() image.Rectangle
	AspectRatio() image.Rectangle
}

type biomAsset struct {
	images      [15]image.Image
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
func NewBiomAsset(srcImages [6]image.Image) (BiomAsset, error) {
	var res image.Rectangle
	for i, img := range srcImages {
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
	dstImages := [15]image.Image{
		// this is old version with 5 images
		// before adding as first image full bottom
		// it uses fewer images to define tile
		// but it flipped horizontal and vertical axis
		// gtexture.NewImage(srcImages[4]).RotateClockwise(3).Image(), // 1000
		// gtexture.NewImage(srcImages[4]).RotateClockwise(0).Image(), // 0100
		// gtexture.NewImage(srcImages[2]).RotateClockwise(1).Image(), // 1100
		// gtexture.NewImage(srcImages[4]).RotateClockwise(2).Image(), // 0010
		// gtexture.NewImage(srcImages[2]).RotateClockwise(0).Image(), // 1010
		// gtexture.NewImage(srcImages[3]).RotateClockwise(2).Image(), // 0110
		// gtexture.NewImage(srcImages[1]).RotateClockwise(1).Image(), // 1110
		// gtexture.NewImage(srcImages[4]).RotateClockwise(1).Image(), // 0001
		// gtexture.NewImage(srcImages[3]).RotateClockwise(1).Image(), // 1001
		// gtexture.NewImage(srcImages[2]).RotateClockwise(2).Image(), // 0101
		// gtexture.NewImage(srcImages[1]).RotateClockwise(2).Image(), // 1101
		// gtexture.NewImage(srcImages[2]).RotateClockwise(3).Image(), // 0011
		// gtexture.NewImage(srcImages[1]).RotateClockwise(0).Image(), // 1011
		// gtexture.NewImage(srcImages[1]).RotateClockwise(3).Image(), // 0111
		// gtexture.NewImage(srcImages[0]).RotateClockwise(0).Image(), // 1111

		// this is new version of tiles
		// it uses one more tile to define biom
		// but it maintains horizontal and vertical axis
		gtexture.NewImage(srcImages[5]).FlipH().Image(),  // 1000
		gtexture.NewImage(srcImages[5]).Image(),          // 0100
		gtexture.NewImage(srcImages[0]).Image(),          // 1100
		gtexture.NewImage(srcImages[5]).FlipHV().Image(), // 0010
		gtexture.NewImage(srcImages[3]).Image(),          // 1010
		gtexture.NewImage(srcImages[4]).Image(),          // 0110
		gtexture.NewImage(srcImages[2]).FlipV().Image(),  // 1110
		gtexture.NewImage(srcImages[5]).FlipV().Image(),  // 0001
		gtexture.NewImage(srcImages[4]).FlipH().Image(),  // 1001
		gtexture.NewImage(srcImages[3]).FlipH().Image(),  // 0101
		gtexture.NewImage(srcImages[2]).FlipHV().Image(), // 1101
		gtexture.NewImage(srcImages[0]).FlipV().Image(),  // 0011
		gtexture.NewImage(srcImages[2]).Image(),          // 1011
		gtexture.NewImage(srcImages[2]).FlipH().Image(),  // 0111
		gtexture.NewImage(srcImages[1]).Image(),          // 1111
	}

	// add this to test:
	// b := map[bool]int{false: 0, true: 1}
	// s := &strings.Builder{}
	// for i, img := range dstImages {
	// 	bounds := img.Bounds()
	// 	_, _, _, lt := img.At(bounds.Min.X, bounds.Min.Y).RGBA()
	// 	_, _, _, rt := img.At(bounds.Max.X-1, bounds.Min.Y).RGBA()
	// 	_, _, _, lb := img.At(bounds.Min.X, bounds.Max.Y-1).RGBA()
	// 	_, _, _, rb := img.At(bounds.Max.X-1, bounds.Max.Y-1).RGBA()
	// 	fmt.Fprintf(s, "i: %2d, lt rt lb rb %d%d%d%d\n", i, b[lt > 0], b[rt > 0], b[lb > 0], b[rb > 0])
	// }
	// fmt.Fprintf(s, "\n\n")
	// print(s.String())

	asset := &biomAsset{
		images:      dstImages,
		res:         res,
		aspectRatio: aspectRatio,
	}
	return asset, nil
}

func (a *biomAsset) Images() [15]image.Image      { return a.images }
func (a *biomAsset) Res() image.Rectangle         { return a.res }
func (a *biomAsset) AspectRatio() image.Rectangle { return a.aspectRatio }
func (a *biomAsset) Release()                     {}
