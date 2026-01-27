package test

import (
	"core/modules/tile"
	"image"
	"image/color"
	"testing"
)

func newImage(v uint8) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))
	if v&0b1000 != 0 {
		img.Set(0, 0, color.White)
	}
	if v&0b0100 != 0 {
		img.Set(1, 0, color.White)
	}
	if v&0b0010 != 0 {
		img.Set(0, 1, color.White)
	}
	if v&0b0001 != 0 {
		img.Set(1, 1, color.White)
	}
	return img
}

func TestNewAsset(t *testing.T) {
	srcImages := [6][]image.Image{
		{newImage(0b0011)},
		{newImage(0b1111)},
		{newImage(0b1110)},
		{newImage(0b1010)},
		{newImage(0b1001)},
		{newImage(0b0001)},
	}

	biom, err := tile.NewBiomAsset(srcImages)
	if err != nil {
		t.Error(err)
	}
	images := biom.Images()

	for i, imgs := range images {
		img := imgs[0]
		bounds := img.Bounds()
		// images are flipped so t and b are swapped
		_, _, _, lt := img.At(bounds.Min.X, bounds.Max.Y-1).RGBA()
		_, _, _, rt := img.At(bounds.Max.X-1, bounds.Max.Y-1).RGBA()
		_, _, _, lb := img.At(bounds.Min.X, bounds.Min.Y).RGBA()
		_, _, _, rb := img.At(bounds.Max.X-1, bounds.Min.Y).RGBA()
		v := 0
		if rb > 0 {
			v |= 1 << 3
		}
		if lb > 0 {
			v |= 1 << 2
		}
		if rt > 0 {
			v |= 1 << 1
		}
		if lt > 0 {
			v |= 1 << 0
		}
		expected := i + 1
		if v != expected {
			t.Errorf("expected %4b but got %4b", expected, v)
		}
	}
}
