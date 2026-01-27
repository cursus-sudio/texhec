package gtexture

import (
	"image"
	"image/draw"
)

type Image interface {
	Image() image.Image

	FlipH() Image
	FlipV() Image
	// horizontally and vertically
	FlipHV() Image

	// rotates 90 deg clockwise
	RotateClockwise(times int) Image

	TrimTransparentBackground() Image
}

type img struct {
	img image.Image
}

func NewImage(image image.Image) Image {
	return &img{img: image}
}

func (s *img) Image() image.Image {
	return s.img
}

func (s *img) FlipH() Image {
    bounds := s.img.Bounds()
    newImg := image.NewRGBA(bounds)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            // Corrected coordinate mapping
            oldX := bounds.Max.X + bounds.Min.X - x - 1
            newImg.Set(x, y, s.img.At(oldX, y))
        }
    }
    s.img = newImg
    return s
}

func (s *img) FlipV() Image {
    bounds := s.img.Bounds()
    newImg := image.NewRGBA(bounds)

    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        // Corrected coordinate mapping
        oldY := bounds.Max.Y + bounds.Min.Y - y - 1
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            newImg.Set(x, y, s.img.At(x, oldY))
        }
    }
    s.img = newImg
    return s
}

func (s *img) FlipHV() Image {
	return s.FlipH().FlipV()
}

func (s *img) RotateClockwise(times int) Image {
	for range times % 4 {
		bounds := s.img.Bounds()
		newBounds := image.Rect(0, 0, bounds.Dy(), bounds.Dx())
		newImg := image.NewRGBA(newBounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				newImg.Set(bounds.Max.Y-y-1, x, s.img.At(x, y))
			}
		}
		s.img = newImg
	}
	return s
}

func (s *img) TrimTransparentBackground() Image {
	bounds := s.img.Bounds()
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := s.img.At(x, y).RGBA()

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
		return s
	}

	newBounds := image.Rect(0, 0, maxX-minX+1, maxY-minY+1)
	croppedImg := image.NewRGBA(newBounds)

	sourcePoint := image.Point{minX, minY}
	draw.Draw(croppedImg, croppedImg.Bounds(), s.img, sourcePoint, draw.Src)

	s.img = croppedImg
	return s
}
