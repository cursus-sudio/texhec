package gameassets

import (
	"image"
	"image/draw"
)

func Rotate90Clockwise(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	W := bounds.Dx()
	H := bounds.Dy()

	rotatedRect := image.Rect(0, 0, H, W)
	rotatedImg := image.NewRGBA(rotatedRect)

	for x := 0; x < W; x++ {
		for y := 0; y < H; y++ {
			c := img.At(x, y)

			newX := H - 1 - y
			newY := x

			rotatedImg.Set(newX, newY, c)
		}
	}

	return rotatedImg
}

func TrimTransparentBackground(img image.Image) image.Image {
	bounds := img.Bounds()
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()

			if a > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	if minX > maxX || minY > maxY { // image is fully transparent there is nothing to trim
		return img
	}

	newBounds := image.Rect(0, 0, maxX-minX+1, maxY-minY+1)
	croppedImg := image.NewRGBA(newBounds)

	sourcePoint := image.Point{minX, minY}
	draw.Draw(croppedImg, croppedImg.Bounds(), img, sourcePoint, draw.Src)

	return croppedImg
}
