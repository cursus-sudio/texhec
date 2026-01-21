package tile

import (
	"core/modules/definition"
	"engine/modules/render"
	"engine/services/assets"
	"engine/services/datastructures"
	gtexture "engine/services/graphics/texture"
	"image"
)

type TileAssets interface {
	AddType(addedAssets datastructures.SparseArray[definition.DefinitionID, assets.AssetID])
}

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
func NewBiomAsset(srcImages [5]image.Image) (BiomAsset, error) {
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
		gtexture.RotateClockwise(srcImages[4], 2), // 1000
		gtexture.RotateClockwise(srcImages[4], 3), // 0100
		gtexture.RotateClockwise(srcImages[2], 1), // 1100
		gtexture.RotateClockwise(srcImages[4], 1), // 0010
		gtexture.RotateClockwise(srcImages[2], 0), // 1010
		gtexture.RotateClockwise(srcImages[3], 1), // 0110
		gtexture.RotateClockwise(srcImages[1], 0), // 1110
		gtexture.RotateClockwise(srcImages[4], 0), // 0001
		gtexture.RotateClockwise(srcImages[3], 0), // 1001
		gtexture.RotateClockwise(srcImages[2], 2), // 0101
		gtexture.RotateClockwise(srcImages[1], 1), // 1101
		gtexture.RotateClockwise(srcImages[2], 3), // 0011
		gtexture.RotateClockwise(srcImages[1], 3), // 1011
		gtexture.RotateClockwise(srcImages[1], 2), // 0111
		gtexture.RotateClockwise(srcImages[0], 0), // 1111
	}

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
