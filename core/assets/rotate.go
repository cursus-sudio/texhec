package gameassets

import "image"

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
