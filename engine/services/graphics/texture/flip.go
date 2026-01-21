package gtexture

import (
	"image"
	"image/draw"
)

func FlipImage(img image.Image) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := 0; y < bounds.Dy(); y++ {
		destY := bounds.Dy() - 1 - y
		destRect := image.Rect(0, destY, bounds.Dx(), destY+1)
		srcRect := image.Rect(0, y, bounds.Dx(), y+1)
		draw.Draw(newImg, destRect, img, srcRect.Min, draw.Src)
	}
	return newImg
}
