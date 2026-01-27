package gtexture

import "image"

func RotateClockwise(img image.Image, times int) image.Image {
	times %= 4
	if times < 0 {
		times += 4
	}
	if times == 0 {
		return img
	}

	b := img.Bounds()
	var newW, newH int

	// Determine final dimensions
	if times%2 == 0 {
		newW, newH = b.Dx(), b.Dy()
	} else {
		newW, newH = b.Dy(), b.Dx()
	}

	res := image.NewRGBA(image.Rect(0, 0, newW, newH))

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			var targetX, targetY int

			nx, ny := x-b.Min.X, y-b.Min.Y

			switch times {
			case 1:
				targetX, targetY = (b.Dy()-1)-ny, nx
			case 2:
				targetX, targetY = (b.Dx()-1)-nx, (b.Dy()-1)-ny
			case 3:
				targetX, targetY = ny, (b.Dx()-1)-nx
			}

			res.Set(targetX, targetY, img.At(x, y))
		}
	}
	return res
}
