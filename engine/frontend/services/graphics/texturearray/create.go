package texturearray

import (
	"image"
	"image/draw"

	"github.com/go-gl/gl/v4.5-core/gl"
)

func createTexs(w, h int, imgs []image.Image) uint32 {
	var texs uint32

	gl.GenTextures(1, &texs)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, texs)
	gl.TexStorage3D(gl.TEXTURE_2D_ARRAY, 1, gl.RGBA8, int32(w), int32(h), int32(len(imgs)))

	for i, img := range imgs {
		rgbaImg, ok := img.(*image.RGBA)
		if !ok {
			rgbaImg = image.NewRGBA(img.Bounds())
			draw.Draw(rgbaImg, rgbaImg.Bounds(), img, image.Point{}, draw.Src)
		}

		gl.TexSubImage3D(gl.TEXTURE_2D_ARRAY, 0, 0, 0, int32(i), int32(w), int32(h), 1, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgbaImg.Pix))
	}

	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D_ARRAY, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.BindTexture(gl.TEXTURE_2D_ARRAY, 0)

	return texs
}
